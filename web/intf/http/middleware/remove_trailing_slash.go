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

package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// RemoveTrailingSlash will remove a trailing slash.
func RemoveTrailingSlash() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		originalURL := utils.CopyString(ctx.OriginalURL())
		if len(originalURL) > 1 && strings.HasSuffix(originalURL, "/") {
			return ctx.Redirect(strings.TrimSuffix(originalURL, "/"))
		}
		return ctx.Next()
	}
}
