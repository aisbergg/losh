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
