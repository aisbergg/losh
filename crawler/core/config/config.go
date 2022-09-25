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

package config

import (
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"
)

type Config struct {
	Crawler  CrawlerConfig `json:"crawler"`
	Log      log.Config    `json:"log"`
	Database dgraph.Config `json:"database"`
}

func DefaultConfig() Config {
	return Config{
		Crawler:  DefaultCrawlerConfig(),
		Log:      log.DefaultConfig(),
		Database: dgraph.DefaultConfig(),
	}
}

type CrawlerConfig struct {
	UserAgent string `json:"userAgent" filter:"trim" validate:"required"`
}

func DefaultCrawlerConfig() CrawlerConfig {
	return CrawlerConfig{
		UserAgent: "LOSH Bot (github.com/aisbergg/losh)",
	}
}
