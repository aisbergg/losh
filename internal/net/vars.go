package net

import "net/http"

var (
	HdrAcceptKey              = http.CanonicalHeaderKey("Accept")
	HdrAuthorizationKey       = http.CanonicalHeaderKey("Authorization")
	HdrContentEncodingKey     = http.CanonicalHeaderKey("Content-Encoding")
	HdrContentLengthKey       = http.CanonicalHeaderKey("Content-Length")
	HdrContentTypeKey         = http.CanonicalHeaderKey("Content-Type")
	HdrLocationKey            = http.CanonicalHeaderKey("Location")
	HdrUserAgentKey           = http.CanonicalHeaderKey("User-Agent")
	HdrXRateLimitRemainingKey = http.CanonicalHeaderKey("X-RateLimit-Remaining")
	HdrXRateLimitResetKey     = http.CanonicalHeaderKey("X-RateLimit-Reset")
)
