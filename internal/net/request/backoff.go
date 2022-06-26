package request

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"

	lerrors "losh/crawler/errors"
)

// ----------------------------------------------------------------------------

// retryable is an interface for a retryable error.
type retryable interface {
	IsRetryable()
	Unwrap() error
}

type retryableError struct {
	err error
}

// NewRetryableError creates a new retryable error.
func NewRetryableError(err error) error {
	return &retryableError{err: err}
}

func (e *retryableError) Error() string {
	return e.err.Error()
}

func (e *retryableError) Unwrap() error {
	return e.err
}

func (e *retryableError) IsRetryable() {}

// httpRetryableError is an error that holds additional response information.
type httpRetryableError struct {
	err  error
	resp *http.Response
}

// NewHTTPRetryableError creates a new http retryable error.
func NewHTTPRetryableError(err error, resp *http.Response) error {
	return &httpRetryableError{err: err, resp: resp}
}

func (e *httpRetryableError) Unwrap() error {
	return e.err
}

func (e *httpRetryableError) Error() string {
	return e.err.Error()
}

func (e *httpRetryableError) IsRetryable() {}

// WithRetryable wraps a backoff function and adds a check for a retryable
// error. When a non retryable error ocurred then no more retry is performed.
func WithRetryable(next retry.Backoff) retry.Backoff {
	return retry.BackoffFunc(func(err error) (time.Duration, error) {
		var rerr retryable
		if !errors.As(err, &rerr) {
			return retry.Stop, err
		}
		return next.Next(rerr.Unwrap())
	})
}

// WithRetryableHTTPResponse works like WithRetryable but also handles http
// errors.
func WithRetryableHTTPResponse(calcDelay CalcWaitTimeFromHTTPResponseFunc, next retry.Backoff) retry.Backoff {
	return retry.BackoffFunc(func(err error) (time.Duration, error) {
		// handle http errors
		var herr *httpRetryableError
		if errors.As(err, &herr) {
			err = herr.Unwrap()

			// get default backoff delay for current attempt
			delay, err := next.Next(err)
			if retry.IsStopped(delay) {
				return retry.Stop, err
			}

			// handle backoff with extra information from response
			delay = calcDelay(herr.resp, delay)

			return delay, err
		}

		// handle other retryable errors
		var rerr retryable
		if errors.As(err, &rerr) {
			return next.Next(rerr.Unwrap())
		}

		// other errors are not retryable
		return retry.Stop, err
	})
}

// WithRetryableHTTPCodes works like WithRetryable but it uses a list of HTTP
// codes to determine if the error is retryable.
func WithRetryableHTTPCodes(retryOnCodes []int, next retry.Backoff) retry.Backoff {
	return retry.BackoffFunc(func(err error) (time.Duration, error) {
		// handle http errors
		var herr *httpRetryableError
		if errors.As(err, &herr) {
			err = herr.Unwrap()

			// check if the response code is in the list of retryable codes
			isRetryableCode := false
			for _, code := range retryOnCodes {
				if herr.resp.StatusCode == code {
					isRetryableCode = true
					break
				}
			}
			if !isRetryableCode {
				return retry.Stop, err
			}

			return next.Next(err)
		}

		// handle other retryable errors
		var rerr retryable
		if errors.As(err, &rerr) {
			return next.Next(rerr.Unwrap())
		}

		// other errors are not retryable
		return retry.Stop, err
	})
}

// WithDelayLimit stops the backoff execution with a temporary error if the
// delay limit is exceeded.
func WithDelayLimit(limit time.Duration, next retry.Backoff) retry.Backoff {
	return retry.BackoffFunc(func(err error) (time.Duration, error) {
		delay, err := next.Next(err)
		if retry.IsStopped(delay) {
			return retry.Stop, err
		}

		if delay > limit {
			retryAfter := time.Now().UTC().Add(delay)
			return retry.Stop, lerrors.NewTemporaryErrorWrap(err, retryAfter, "wait time limit exceeded")
		}

		return delay, err
	})
}

// WithLogging wraps a backoff function and logs the retry attempts.
func WithLogging(log *zap.SugaredLogger, next retry.Backoff) retry.Backoff {
	return retry.BackoffFunc(func(err error) (time.Duration, error) {
		delay, err := next.Next(err)
		if retry.IsStopped(delay) {
			return retry.Stop, err
		}

		log.Debugf("wait for %s and retry request: %s", delay.String(), err)

		return delay, err
	})
}

// RetryAfter returns a delay based on the Retry-After header of the response.
// If the header is not present or cannot be parsed, the default delay is
// returned.
func RetryAfter(resp *http.Response, delay time.Duration) time.Duration {
	retryAfterString := resp.Header.Get("Retry-After")
	if retryAfterString == "" {
		return delay
	}

	// interpret Retry-After as seconds
	seconds, err := strconv.ParseInt(retryAfterString, 10, 64)
	if err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}

	// interpret Retry-After as date
	date, err := time.Parse(time.RFC1123, retryAfterString)
	if err == nil {
		retryAfter := date.Sub(time.Now())
		if retryAfter > 0 {
			return retryAfter
		}
	}

	// return default delay if parsing fails
	return delay
}
