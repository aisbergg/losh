package config

import (
	"losh/internal/infra/dgraph"
	"losh/internal/lib/log"
	"time"
)

type Config struct {
	Debug     DebugConfig     `json:"debug"`
	Log       log.Config      `json:"log"`
	AccessLog AccessLogConfig `json:"accessLog"`
	Server    ServerConfig    `json:"server"`
	Database  dgraph.Config   `json:"database"`
}

func DefaultConfig() Config {
	return Config{
		Debug:     DefaultDebugConfig(),
		Log:       log.DefaultConfig(),
		AccessLog: DefaultAccessLogConfig(),
		Server:    DefaultServerConfig(),
		Database:  dgraph.DefaultConfig(),
	}
}

type DebugConfig struct {
	Enabled bool `json:"enabled"`
	Pprof   bool `json:"pprof"`
	Expvar  bool `json:"expvar"`
}

func DefaultDebugConfig() DebugConfig {
	return DebugConfig{
		Enabled: false,
		Pprof:   false,
		Expvar:  false,
	}
}

type AccessLogConfig struct {
	log.Config
	Enabled bool     `json:"enabled"`
	Fields  []string `json:"fields"`
}

func DefaultAccessLogConfig() AccessLogConfig {
	return AccessLogConfig{
		Config:  log.DefaultConfig(),
		Enabled: false,
		Fields:  []string{},
	}
}

type ServerConfig struct {
	Interface      string            `json:"interface"`
	Port           uint16            `json:"port"`
	TrustedDomains []string          `json:"trustedDomains"`
	Compress       int               `json:"compress"`
	Cache          ServerCacheConfig `json:"cache"`
	TLS            TLSConfig         `json:"tls"`
}

func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Interface:      "",
		Port:           3000,
		TrustedDomains: []string{},
		Compress:       1,
		TLS:            DefaultTLSConfig(),
		Cache:          DefaultServerCacheConfig(),
	}
}

type ServerCacheConfig struct {
	Enabled      bool          `json:"enabled"`
	Expiration   time.Duration `json:"expiration"`
	CacheControl bool          `json:"cacheControl"`
}

func DefaultServerCacheConfig() ServerCacheConfig {
	return ServerCacheConfig{
		Enabled:      false,
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}
}

type TLSConfig struct {
	Enabled     bool   `json:"enabled"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`

	// https://github.com/golang/go/blob/master/src/crypto/tls/cipher_suites.go

	MinVersion               string   `json:"minVersion" filter:"trim" validate:"required|in:VersionTLS10,VersionTLS11,VersionTLS12"`
	CipherSuites             []string `json:"cipherSuites"`     // TODO: would be nice to filter/validate on sub elements filter:"trim,upper" validate:"minLen:1|in:TLS_RSA_WITH_RC4_128_SHA,TLS_RSA_WITH_3DES_EDE_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA,TLS_RSA_WITH_AES_256_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA256,TLS_RSA_WITH_AES_128_GCM_SHA256,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_RSA_WITH_RC4_128_SHA,TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,TLS_AES_128_GCM_SHA256,TLS_AES_256_GCM_SHA384,TLS_CHACHA20_POLY1305_SHA256,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305"
	CurvePreferences         []string `json:"curvePreferences"` // filter:"trim,lower" validate:"minLen:1|in:CurveP256,CurveP384,CurveP521,X25519"
	PreferServerCipherSuites bool     `json:"preferServerCipherSuites"`
}

func DefaultTLSConfig() TLSConfig {
	return TLSConfig{
		Enabled:     false,
		Certificate: "",
		Key:         "",
		MinVersion:  "VersionTLS12",
		CipherSuites: []string{
			// TLS 1.3 Cipher Suites
			"TLS_AES_256_GCM_SHA384",
			"TLS_AES_128_GCM_SHA256",
			"TLS_CHACHA20_POLY1305_SHA256",
			// TLS 1.2 Cipher Suites
			"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
			"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256",
			"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256",
		},
		CurvePreferences: []string{
			"X25519",
			"CurveP521",
			"CurveP384",
		},
		PreferServerCipherSuites: false,
	}
}
