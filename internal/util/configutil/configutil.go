package configutil

import (
	"reflect"
	"time"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"

	configLoader "github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/gookit/validate"
)

func Load(path pathlib.Path, out interface{}) error {
	// [OPTIONAL]
	// ENV
	// include prefix
	// iterate over a schema and map env vars to new structure
	// normalize
	// validate

	// FLAGS
	// iterate over a schema and map env vars to new structure
	// normalize
	// validate

	// [PRIORITY]
	// YAML
	// validate

	// MERGE
	//

	// VALIDATE
	// default values

	// configure config loader
	configLoader.WithOptions(func(opt *configLoader.Options) {
		opt.DecoderConfig.TagName = "json"
		opt.DecoderConfig.WeaklyTypedInput = true
		// custom type decoders
		opt.DecoderConfig.DecodeHook = customDecodeHook

	})
	configLoader.AddDriver(yamlv3.Driver)

	// now load config from file
	if path.String() != "." {
		if exists, err := path.Exists(); err != nil || !exists {
			return errors.Errorf("configuration file doesn't exist: %s", path.String())
		}
		if isFile, err := path.IsFile(); err != nil || !isFile {
			return errors.Errorf("given path is not a file: %s", path.String())
		}
		if err := configLoader.LoadFiles(path.String()); err != nil {
			return err
		}
	}

	// bind loaded config to model
	err := configLoader.BindStruct("", out)
	if err != nil {
		return err
	}

	// validate model
	validate.AddGlobalMessages(map[string]string{
		"enum": "value {field} must be one of %v",
	})
	validator := validate.Struct(out)
	if !validator.Validate() {
		return errors.Errorf("configuration is invalid: %s", validator.Errors.Error())
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
