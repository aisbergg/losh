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

	losherrors "losh/internal/lib/errors"
	"losh/internal/lib/net/ratelimit"
)

// Default values for retry parameters.
const (
	DefaultMaxWaitTime        = 5 * time.Minute
	DefaultRetryCount         = 5
	DefaultRetryWaitTime      = 1 * time.Second
	DefaultRetryJitterPercent = 5
)

// Default values for retry parameters.
var (
	DefaultRetryOnCodes = []int{408, 429, 500, 502, 503, 504}
)

// applyRateLimit applies the rate limit to the request.
func applyRateLimiters(ctx context.Context, rateLimiters []ratelimit.RateLimiter, maxWaitTime time.Duration) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	wait := time.Duration(0)
	for _, rl := range rateLimiters {
		rlWait := rl.WaitTime()
		if rlWait > wait {
			wait = rlWait
		}
	}
	if maxWaitTime > 0 && wait > maxWaitTime {
		return losherrors.NewTemporaryError(nil, time.Now().Add(wait), "wait time limit exceeded")
	}
	ratelimit.SleepWithContext(ctx, wait)
	return nil
}
