package request

import (
	"context"
	gourl "net/url"
	"time"

	"losh/internal/lib/net/ratelimit"

	"github.com/aisbergg/go-retry/pkg/retry"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type RESTRequest struct {
	Req    *resty.Request
	Method string
	URL    string
}

func NewRESTRequest(req *resty.Request, method, url string) RESTRequest {
	return RESTRequest{
		Req:    req,
		Method: method,
		URL:    url,
	}
}

type RESTRequester struct {
	// Optional logger for logging failed requests.
	logger *zap.SugaredLogger
	// Underlying REST client.
	restClient *resty.Client

	//
	//  Rate Limit
	//

	// The rate limiters to apply to the requests.
	rateLimiters []ratelimit.RateLimiter
	// Max wait time between requests before a temporary error is returned.
	maxWaitTime time.Duration

	//
	// Retry
	//

	// Max number of retries for a request.
	retryCount uint64
	// Initial wait time between retries.
	retryWaitTime time.Duration
	// Percentage of jitter to add to the retry wait time.
	retryJitterPercent uint64
	// Amount of jitter to add to the retry wait time. `retryJitterPercent`
	// takes precedence.
	retryJitter time.Duration
	// HTTP return codes to retry on. `retryCalcWaitTimeFromResponseFunc` takes
	// precedence.
	retryOnCodes []int
	// Function to calculate the next retry wait time based on the HTTP
	// response.
	retryCalcWaitTimeFromResponseFunc CalcWaitTimeFromHTTPResponseFunc
}

// NewRESTRequester creates a new Requester with given HTTP client.
func NewRESTRequester(client *resty.Client) *RESTRequester {
	return &RESTRequester{
		restClient: client,

		// Rate Limit
		rateLimiters: []ratelimit.RateLimiter{},
		maxWaitTime:  DefaultMaxWaitTime,

		// Retry
		retryCount:         DefaultRetryCount,
		retryWaitTime:      DefaultRetryWaitTime,
		retryJitterPercent: DefaultRetryJitterPercent,
		retryOnCodes:       DefaultRetryOnCodes,
	}
}

// SetLogger sets the logger to use for logging failed requests.
func (r *RESTRequester) SetLogger(logger *zap.SugaredLogger) *RESTRequester {
	r.logger = logger
	return r
}

// AddRateLimiter adds the rate limiter to apply to the requests.
func (r *RESTRequester) AddRateLimiter(rl ratelimit.RateLimiter) *RESTRequester {
	r.rateLimiters = append(r.rateLimiters, rl)
	return r
}

// AddRateLimiters adds the rate limiters to apply to the requests.
func (r *RESTRequester) AddRateLimiters(rls []ratelimit.RateLimiter) *RESTRequester {
	for _, rl := range rls {
		r.rateLimiters = append(r.rateLimiters, rl)
	}
	return r
}

// SetMaxWaitTime method sets the max wait time between requests before a
// temporary error is returned (indicator for "retry later").
//
// Default is 5 minutes.
func (r *RESTRequester) SetMaxWaitTime(maxWaitTime time.Duration) *RESTRequester {
	r.maxWaitTime = maxWaitTime
	return r
}

// SetRetryCount adds the max number of retries to perform for a request.
//
// Default is 5.
func (r *RESTRequester) SetRetryCount(count uint64) *RESTRequester {
	r.retryCount = count
	return r
}

// SetRetryWaitTime sets the initial wait time to sleep before retrying the
// request.
//
// Default is 1 second.
func (r *RESTRequester) SetRetryWaitTime(waitTime time.Duration) *RESTRequester {
	r.retryWaitTime = waitTime
	return r
}

// SetRetryJitterPercent sets the percentage of jitter to add to the retry wait
// time.
//
// Default is 5.
func (r *RESTRequester) SetRetryJitterPercent(percent uint64) *RESTRequester {
	r.retryJitterPercent = percent
	return r
}

// SetRetryJitter sets the jitter to add to the retry wait time.
func (r *RESTRequester) SetRetryJitter(jitter time.Duration) *RESTRequester {
	r.retryJitter = jitter
	return r
}

// SetRetryOnCodes sets the HTTP return codes to retry on.
//
// Default are 408, 429, 500, 502, 503, 504.
func (r *RESTRequester) SetRetryOnCodes(codes []int) *RESTRequester {
	r.retryOnCodes = codes
	return r
}

// SetRetryCalcWaitTimeFromResponseFunc sets the function to calculate the
// next retry wait time based on the HTTP response.
func (r *RESTRequester) SetRetryCalcWaitTimeFromResponseFunc(f CalcWaitTimeFromHTTPResponseFunc) *RESTRequester {
	r.retryCalcWaitTimeFromResponseFunc = f
	return r
}

// Do executes the request and returns the response.
func (r *RESTRequester) Do(req RESTRequest) (resp *resty.Response, err error) {
	ctx := req.Req.Context()

	// apply rate limit beforehand
	if err = applyRateLimiters(ctx, r.rateLimiters, r.maxWaitTime); err != nil {
		return
	}

	// create wrapper for request execution
	retryFunc := func(ctx context.Context) (err error) {
		// check if the request was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// execute request
		resp, err = req.Req.Execute(req.Method, req.URL)

		// update rate limiters after each request
		for _, rl := range r.rateLimiters {
			rl.Update(resp)
		}

		// on error
		if err != nil {
			// on network error -> retry
			if _, ok := err.(*gourl.Error); ok {
				return &retryableError{err: err}
			}

			// on HTTP error -> retry
			if resp != nil {
				for _, c := range r.retryOnCodes {
					if c == resp.RawResponse.StatusCode {
						return &httpRetryableError{
							err:  err,
							resp: resp.RawResponse,
						}
					}
				}
			}
		}

		// stop on success
		return nil
	}

	// execute the request with retries
	err = retry.Do(ctx, r.createBackoff(), retryFunc)

	return
}

// createBackoff creates a new backoff.Backoff instance.
func (r *RESTRequester) createBackoff() retry.Backoff {
	b := retry.NewExponential(r.retryWaitTime)
	if r.retryCalcWaitTimeFromResponseFunc != nil {
		b = WithRetryableHTTPResponse(r.retryCalcWaitTimeFromResponseFunc, b)
	} else if len(r.retryOnCodes) > 0 {
		b = WithRetryableHTTPCodes(r.retryOnCodes, b)
	} else {
		b = WithRetryable(b)
	}
	if r.retryJitterPercent > 0 && r.retryJitterPercent <= 100 {
		b = retry.WithJitterPercent(r.retryJitterPercent, true, b)
	} else if r.retryJitter > 0 {
		b = retry.WithJitter(r.retryJitter, true, b)
	}
	if r.maxWaitTime > 0 {
		b = WithDelayLimit(r.maxWaitTime, b)
	}
	if r.retryCount > 0 {
		b = retry.WithMaxRetries(r.retryCount, b)
	}
	if r.logger != nil {
		b = WithLogging(r.logger, b)
	}
	return b
}
