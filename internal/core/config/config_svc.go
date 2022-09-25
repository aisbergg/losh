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
	"io"
	"os"
	"reflect"
	"time"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"

	configLoader "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/validate"
	"gopkg.in/yaml.v3"
)

// Service provides methods to load, store and validate the configuration.
type Service struct {
	path     pathlib.Path
	defCfgFn func() interface{}
}

// NewService creates a new configuration service.
func NewService(path string, defCfgFn func() interface{}) Service {
	return Service{
		path:     pathlib.NewPath(path),
		defCfgFn: defCfgFn,
	}
}

// Init initializes a default configuration.
func (s Service) Init() error {
	return s.Save(s.GetDefault())
}

// Validate validates the configuration. If the configuration is invalid it
// returns an error with details.
func (Service) Validate(cfg interface{}) error {
	validate.AddGlobalMessages(map[string]string{
		"enum": "value {field} must be one of %v",
	})
	validator := validate.Struct(cfg)
	if !validator.Validate() {
		return errors.Errorf("configuration is invalid: %s", validator.Errors.Error())
	}
	return nil
}

// GetDefault returns the default configuration.
func (s Service) GetDefault() interface{} {
	return s.defCfgFn()
}

// Get loads, validates and returns the configuration.
func (s Service) Get() (interface{}, error) {
	// default configuration
	cfg := s.defCfgFn()

	// load from file
	if err := s.getFromFile(cfg); err != nil {
		return nil, err
	}

	// validate model
	if err := s.Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getFromFile loads the configuration from the configured file.
func (s Service) getFromFile(bind interface{}) error {
	// configure config loader
	configLoader.WithOptions(func(opt *configLoader.Options) {
		opt.DecoderConfig.TagName = "json"
		opt.DecoderConfig.WeaklyTypedInput = true
		// custom type decoders
		opt.DecoderConfig.DecodeHook = customDecodeHook

	})
	configLoader.AddDriver(yamlv3.Driver)

	// now load config from file
	if s.path.String() != "." {
		if exists, err := s.path.Exists(); err != nil || !exists {
			return errors.Errorf("configuration file doesn't exist: %s", s.path.String())
		}
		if isFile, err := s.path.IsFile(); err != nil || !isFile {
			return errors.Errorf("given path is not a file: %s", s.path.String())
		}
		if err := configLoader.LoadFiles(s.path.String()); err != nil {
			return errors.Wrapf(err, "failed to load file: %s", s.path.String())
		}
	}

	// bind loaded config to model
	err := configLoader.BindStruct("", bind)
	if err != nil {
		return err
	}

	return nil
}

// Save saves the configuration to the configured file.
func (s Service) Save(cfg interface{}) error {
	// create file to write to
	var writeTo io.Writer
	if s.path.String() != "." {
		exists, err := s.path.Exists()
		if err != nil {
			return err
		}
		if exists {
			isFile, err := s.path.IsFile()
			if err != nil {
				return err
			}
			if !isFile {
				return errors.Errorf("path is not a file: %s", s.path.String())
			}
		}
		file, err := s.path.OpenFile(os.O_RDWR | os.O_CREATE)
		if err != nil {
			return err
		}
		defer file.Close()
		writeTo = file

	} else {
		// write to stdout
		writeTo = os.Stdout
	}

	// encode as yaml
	yamlEncoder := yaml.NewEncoder(writeTo)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}

func customDecodeHook(from reflect.Value, to reflect.Value) (interface{}, error) {
	switch to.Type() {
	case reflect.TypeOf(time.Duration(0)):
		if condition := from.Kind() == reflect.String; condition {
			return time.ParseDuration(from.String())
		}
	}
	return from.Interface(), nil
}
