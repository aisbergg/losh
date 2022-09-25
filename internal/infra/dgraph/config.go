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

package dgraph

import "time"

type Config struct {
	Host    string          `json:"host"`
	Port    uint64          `json:"port"`
	TLS     ClientTLSConfig `json:"tls"`
	Timeout time.Duration   `json:"timeout"`
}

func DefaultConfig() Config {
	return Config{
		Host:    "localhost",
		Port:    8080,
		TLS:     DefaultClientTLSConfig(),
		Timeout: 60 * time.Second,
	}
}

type ClientTLSConfig struct {
	Enabled bool `json:"enabled"`
	// Verify the server's certificate against the list of supplied CAs.
	Verify      bool   `json:"verify"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
	// CaCertificates is a list of trusted root certificate authorities for verifying the servers certificate.
	CACertificates []string `json:"caCertificates"`
}

func DefaultClientTLSConfig() ClientTLSConfig {
	return ClientTLSConfig{
		Enabled:        false,
		Verify:         false,
		Certificate:    "",
		Key:            "",
		CACertificates: []string{},
	}
}
