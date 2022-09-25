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

package ratelimit

import (
	"context"
	"math/rand"
	"time"
)

// RateLimiter applies a rate limit by waiting if necessary.
type RateLimiter interface {
	// Update updates the rate limit after each request. It can make use of
	// information from the response to calculate the next time limit.
	Update(resp interface{})

	// WaitTime returns the amount of time to wait before the next request.
	WaitTime() time.Duration

	// Apply applies the rate limit by waiting if necessary. The context is used
	// to cancel the wait.
	Apply(context.Context)
}

// -----------------------------------------------------------------------------
// RateLimitTimedelta
// -----------------------------------------------------------------------------

// TimedeltaRateLimiter is a `RateLimiter` that ensures that requests are at least
// `timedelta` time apart from each other. It adds a random amount of jitter to
// the wait time.
type TimedeltaRateLimiter struct {
	lastTime  time.Time
	timedelta time.Duration
	maxJitter time.Duration
}

// NewTimedeltaRateLimiter creates a new `TimedeltaRateLimiter`.
func NewTimedeltaRateLimiter(timedelta time.Duration, jitterPercent uint64) *TimedeltaRateLimiter {
	maxJitter := time.Duration(0)
	if jitterPercent > 0 {
		maxJitter = time.Duration(float64(timedelta) * float64(jitterPercent) / 100.0)
	}

	return &TimedeltaRateLimiter{
		lastTime:  time.Time{},
		timedelta: timedelta,
		maxJitter: maxJitter,
	}
}

// Update implements `RateLimiter`.
func (rl *TimedeltaRateLimiter) Update(_ interface{}) {
	rl.lastTime = time.Now().UTC()
}

// WaitTime implements `RateLimiter`.
func (rl *TimedeltaRateLimiter) WaitTime() time.Duration {
	jitter := time.Duration(0)
	if rl.maxJitter > 0 {
		jitter = time.Duration(rand.Int63n(int64(rl.maxJitter)))
	}

	wait := rl.timedelta + jitter - time.Now().UTC().Sub(rl.lastTime)
	if wait > 0 {
		return wait
	}
	return 0
}

// Apply implements RateLimiter.
func (rl *TimedeltaRateLimiter) Apply(ctx context.Context) {
	SleepWithContext(ctx, rl.WaitTime())
}

// -----------------------------------------------------------------------------
// NumRequestsRateLimiter
// -----------------------------------------------------------------------------

// NumRequestsRateLimiter is a `RateLimiter` that limits the number of requests.
type NumRequestsRateLimiter struct {
	numRequests       uint64
	requestsRemaining uint64
	waitTime          time.Duration
}

// NewNumRequestsRateLimiter creates a new `NumRequestsRateLimiter`.
func NewNumRequestsRateLimiter(numRequests uint64, waitTime time.Duration) *NumRequestsRateLimiter {
	return &NumRequestsRateLimiter{
		numRequests:       numRequests,
		requestsRemaining: numRequests,
		waitTime:          waitTime,
	}
}

// Update implements `RateLimiter`.
func (rl *NumRequestsRateLimiter) Update(resp interface{}) {
	if rl.requestsRemaining > 0 {
		rl.requestsRemaining--
	} else {
		// reset the counter
		rl.requestsRemaining = rl.numRequests
	}
}

// WaitTime implements `RateLimiter`.
func (rl *NumRequestsRateLimiter) WaitTime() time.Duration {
	if rl.requestsRemaining > 0 {
		return 0
	}
	return rl.waitTime
}

// Apply implements `RateLimiter`.
func (rl *NumRequestsRateLimiter) Apply(ctx context.Context) {
	SleepWithContext(ctx, rl.WaitTime())
}

// -----------------------------------------------------------------------------

// SleepWithContext sleeps for the given duration, but will return early if the
// context is canceled.
func SleepWithContext(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}
