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
		UserAgent: "OKH-LOSH-Crawler github.com/OPEN-NEXT/OKH-LOSH",
	}
}
