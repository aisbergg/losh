package config

import (
	"losh/internal/logging"
	"losh/internal/repository/dgraph"
)

type Config struct {
	Crawler  CrawlerConfig  `json:"crawler"`
	Log      logging.Config `json:"log"`
	Database dgraph.Config  `json:"database"`
}

func DefaultConfig() Config {
	return Config{
		Crawler:  DefaultCrawlerConfig(),
		Log:      logging.DefaultConfig(),
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
