package config

import "losh/internal/core/config"

// Service provides methods to load, store and validate the configuration.
type Service struct {
	config.Service
}

// NewService creates a new configuration service.
func NewService(path string) Service {
	defCfgFn := func() interface{} { return DefaultConfig() }
	return Service{
		Service: config.NewService(path, defCfgFn),
	}
}

// Get loads, validates and returns the configuration.
func (s Service) Get() (Config, error) {
	cfg, err := s.Service.Get()
	if err != nil {
		return Config{}, err
	}
	return cfg.(Config), nil
}

// Save saves the configuration to the configured file.
func (s Service) Save(cfg Config) error {
	return s.Service.Save(cfg)
}
