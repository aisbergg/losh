// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package request

import (
	"context"
	"net/http"
	"time"

	"losh/internal/lib/net/ratelimit"

	"github.com/aisbergg/go-retry/pkg/retry"
	"go.uber.org/zap"
)

// CalcWaitTimeFromHTTPResponseFunc calculates the delay based on the HTTP
// response. A negative delay indicates that the backoff should not retry the
// request.
type CalcWaitTimeFromHTTPResponseFunc func(resp *http.Response, delay time.Duration) time.Duration

//
// Functionality
//

// - rate limit apply
// - then let backoff take over?
// - backoff includes rate limit apply? - need to update and calculate times

// backoff when
// - error in transmission
// - error code [429, 500, 502, 503, 504]
//   - (Retry-After)
//   - (fixed time) (min time)
// - error code 403 and "rate limit" (60 seconds)
// self._primary_search_rate_limit.update(
// 	num_requests=int(response.headers["X-RateLimit-Remaining"]),
// 	reset_time=datetime.fromtimestamp(int(response.headers["X-RateLimit-Reset"]), tz=timezone.utc),
// )

type HTTPRequest struct {
	Ctx           context.Context
	OperationName string
	Query         string
	Variables     map[string]interface{}
}

func NewHTTPRequest(ctx context.Context, operationName string, query string, vars map[string]interface{}) HTTPRequest {
	return HTTPRequest{
		Ctx:           ctx,
		OperationName: operationName,
		Query:         query,
		Variables:     vars,
	}
}

type HTTPRequester struct {
	// Optional logger for logging failed requests.
	logger *zap.SugaredLogger
	// Underlying HTTP httpClient.
	httpClient *http.Client

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
	retryCount int
	// Initial wait time between retries.
	retryWaitTime time.Duration
	// Percentage of jitter to add to the retry wait time.
	retryJitterPercent int
	// Amount of jitter to add to the retry wait time.
	retryJitter time.Duration
	// HTTP return codes to retry on. If `retryCalcWaitTimeFromResponseFunc` is
	// set then this is ignored.
	retryOnCodes []int
	// Function to calculate the next retry wait time based on the HTTP
	// response.
	retryCalcWaitTimeFromResponseFunc CalcWaitTimeFromHTTPResponseFunc
}

// NewHTTPRequester creates a new Requester with given HTTP client.
func NewHTTPRequester(client *http.Client) *HTTPRequester {
	return &HTTPRequester{
		httpClient: client,

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
func (r *HTTPRequester) SetLogger(logger *zap.SugaredLogger) *HTTPRequester {
	r.logger = logger
	return r
}

// AddRateLimiter adds the rate limiter to apply to the requests.
func (r *HTTPRequester) AddRateLimiter(rl ratelimit.RateLimiter) *HTTPRequester {
	r.rateLimiters = append(r.rateLimiters, rl)
	return r
}

// AddRateLimiters adds the rate limiters to apply to the requests.
func (r *HTTPRequester) AddRateLimiters(rls []ratelimit.RateLimiter) *HTTPRequester {
	for _, rl := range rls {
		r.rateLimiters = append(r.rateLimiters, rl)
	}
	return r
}

// SetMaxWaitTime method sets the max wait time between requests before a
// temporary error is returned (indicator for "retry later").
//
// Default is 5 minutes.
func (r *HTTPRequester) SetMaxWaitTime(maxWaitTime time.Duration) *HTTPRequester {
	r.maxWaitTime = maxWaitTime
	return r
}

// SetRetryCount adds the max number of retries to perform for a request.
//
// Default is 5.
func (r *HTTPRequester) SetRetryCount(count int) *HTTPRequester {
	r.retryCount = count
	return r
}

// SetRetryWaitTime sets the initial wait time to sleep before retrying the
// request.
//
// Default is 1 second.
func (r *HTTPRequester) SetRetryWaitTime(waitTime time.Duration) *HTTPRequester {
	r.retryWaitTime = waitTime
	return r
}

// SetRetryJitterPercent sets the percentage of jitter to add to the retry wait
// time.
//
// Default is 5.
func (r *HTTPRequester) SetRetryJitterPercent(percent int) *HTTPRequester {
	r.retryJitterPercent = percent
	return r
}

// SetRetryJitter sets the jitter to add to the retry wait time.
func (r *HTTPRequester) SetRetryJitter(jitter time.Duration) *HTTPRequester {
	r.retryJitter = jitter
	return r
}

// SetRetryOnCodes sets the HTTP return codes to retry on.
//
// Default are 408, 429, 500, 502, 503, 504.
func (r *HTTPRequester) SetRetryOnCodes(codes []int) *HTTPRequester {
	r.retryOnCodes = codes
	return r
}

// SetRetryCalcWaitTimeFromResponseFunc sets the function to calculate the
// next retry wait time based on the HTTP response.
func (r *HTTPRequester) SetRetryCalcWaitTimeFromResponseFunc(f CalcWaitTimeFromHTTPResponseFunc) *HTTPRequester {
	r.retryCalcWaitTimeFromResponseFunc = f
	return r
}

// Do performs the HTTP request and returns the response.
func (r *HTTPRequester) Do(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()

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
		resp, err = r.httpClient.Do(req)

		// update rate limiters after each request
		for _, rl := range r.rateLimiters {
			rl.Update(resp)
		}

		// on network error
		if err != nil {
			return &retryableError{err: err}
		}

		// on http request error
		if resp.StatusCode != 200 {
			return &httpRetryableError{
				err:  err,
				resp: resp,
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
func (r *HTTPRequester) createBackoff() retry.Backoff {
	b := retry.NewExponential(r.retryWaitTime)
	if r.retryCalcWaitTimeFromResponseFunc != nil {
		b = WithRetryableHTTPResponse(r.retryCalcWaitTimeFromResponseFunc, b)
	} else if len(r.retryOnCodes) > 0 {
		b = WithRetryableHTTPCodes(r.retryOnCodes, b)
	} else {
		b = WithRetryable(b)
	}
	if r.retryJitterPercent > 0 && r.retryJitterPercent <= 100 {
		b = retry.WithJitterPercent(uint64(r.retryJitterPercent), true, b)
	} else if r.retryJitter > 0 {
		b = retry.WithJitter(r.retryJitter, true, b)
	}
	if r.maxWaitTime > 0 {
		b = WithDelayLimit(r.maxWaitTime, b)
	}
	if r.retryCount > 0 {
		b = retry.WithMaxRetries(uint64(r.retryCount), b)
	}
	if r.logger != nil {
		b = WithLogging(r.logger, b)
	}
	return b
}
