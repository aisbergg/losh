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
	"time"

	"losh/internal/lib/net/ratelimit"

	gql "github.com/Yamashou/gqlgenc/clientv2"
	"github.com/aisbergg/go-retry/pkg/retry"
	"go.uber.org/zap"
)

type GraphQLRequest struct {
	Ctx           context.Context
	OperationName string
	Query         string
	Variables     map[string]interface{}
}

func NewGraphQLRequest(ctx context.Context, operationName string, query string, vars map[string]interface{}) GraphQLRequest {
	return GraphQLRequest{
		Ctx:           ctx,
		OperationName: operationName,
		Query:         query,
		Variables:     vars,
	}
}

type GraphQLRequester struct {
	// Optional logger for logging failed requests.
	logger *zap.SugaredLogger
	// The underlying GraphQL gqlClient.
	gqlClient *gql.Client

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
	// HTTP return codes to retry on.
	retryOnCodes []int
}

// NewGraphQLRequester creates a new Requester with given HTTP client.
func NewGraphQLRequester(client *gql.Client) *GraphQLRequester {
	return &GraphQLRequester{
		gqlClient: client,

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
func (r *GraphQLRequester) SetLogger(logger *zap.SugaredLogger) *GraphQLRequester {
	r.logger = logger
	return r
}

// AddRateLimiter adds the rate limiter to apply to the requests.
func (r *GraphQLRequester) AddRateLimiter(rl ratelimit.RateLimiter) *GraphQLRequester {
	r.rateLimiters = append(r.rateLimiters, rl)
	return r
}

// AddRateLimiters adds the rate limiters to apply to the requests.
func (r *GraphQLRequester) AddRateLimiters(rls []ratelimit.RateLimiter) *GraphQLRequester {
	for _, rl := range rls {
		r.rateLimiters = append(r.rateLimiters, rl)
	}
	return r
}

// SetMaxWaitTime method sets the max wait time between requests before a
// temporary error is returned (indicator for "retry later").
//
// Default is 5 minutes.
func (r *GraphQLRequester) SetMaxWaitTime(maxWaitTime time.Duration) *GraphQLRequester {
	r.maxWaitTime = maxWaitTime
	return r
}

// SetRetryCount adds the max number of retries to perform for a request.
//
// Default is 5.
func (r *GraphQLRequester) SetRetryCount(count uint64) *GraphQLRequester {
	r.retryCount = count
	return r
}

// SetRetryWaitTime sets the initial wait time to sleep before retrying the
// request.
//
// Default is 1 second.
func (r *GraphQLRequester) SetRetryWaitTime(waitTime time.Duration) *GraphQLRequester {
	r.retryWaitTime = waitTime
	return r
}

// SetRetryJitterPercent sets the percentage of jitter to add to the retry wait
// time.
//
// Default is 5.
func (r *GraphQLRequester) SetRetryJitterPercent(percent uint64) *GraphQLRequester {
	r.retryJitterPercent = percent
	return r
}

// SetRetryJitter sets the jitter to add to the retry wait time.
func (r *GraphQLRequester) SetRetryJitter(jitter time.Duration) *GraphQLRequester {
	r.retryJitter = jitter
	return r
}

// SetRetryOnCodes sets the HTTP return codes to retry on.
//
// Default are 408, 429, 500, 502, 503, 504.
func (r *GraphQLRequester) SetRetryOnCodes(codes []int) *GraphQLRequester {
	r.retryOnCodes = codes
	return r
}

// Do executes the request and returns the response.
func (r *GraphQLRequester) Do(req GraphQLRequest, resp interface{}) (err error) {
	// apply rate limit beforehand
	if err = applyRateLimiters(req.Ctx, r.rateLimiters, r.maxWaitTime); err != nil {
		return
	}

	// create wrapper for request execution
	retryFunc := func(ctx context.Context) error {
		// check if the request was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// execute request
		err = r.gqlClient.Post(ctx, req.OperationName, req.Query, resp, req.Variables)

		// update rate limiters after each request
		for _, rl := range r.rateLimiters {
			rl.Update(resp)
		}

		if err != nil {
			// error parsing the response
			if errResponse, ok := err.(*gql.ErrorResponse); ok {
				if errResponse.NetworkError != nil {
					// on HTTP error -> retry
					for _, c := range r.retryOnCodes {
						if c == errResponse.NetworkError.Code {
							return NewRetryableError(errResponse)
						}
					}
				}

				// on untryable HTTP error or graphql error -> stop
				return errResponse
			}

			// on network error -> retry
			return NewRetryableError(err)
		}

		// stop on success
		return nil
	}

	// execute the request with retries
	err = retry.Do(req.Ctx, r.createBackoff(), retryFunc)

	return
}

// createBackoff creates a new backoff.Backoff instance.
func (r *GraphQLRequester) createBackoff() retry.Backoff {
	b := retry.NewExponential(r.retryWaitTime)
	if len(r.retryOnCodes) > 0 {
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
