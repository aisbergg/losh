package dgraph

type Config struct {
	Host string          `json:"host"`
	Port uint64          `json:"port"`
	TLS  ClientTLSConfig `json:"tls"`
}

func DefaultConfig() Config {
	return Config{
		Host: "localhost",
		Port: 8080,
		TLS:  DefaultClientTLSConfig(),
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
