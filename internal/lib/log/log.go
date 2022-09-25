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

package log

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerManager is just for wrapping the logging instance together with the log
// file.
type LoggerManager struct {
	rootLogger *zap.Logger
	logFile    *lumberjack.Logger
}

// NewLoggerManager creates a new LoggerManager with the given configuration.
func NewLoggerManager(config Config) (LoggerManager, error) {
	lm := LoggerManager{}

	var encoder zapcore.Encoder
	switch config.Format {
	case "json":
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.TimeKey = "time"
		encoderConfig.NameKey = "name"
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		panic(fmt.Sprintf("unknown logging format: %s", config.Format))
	}

	var writer zapcore.WriteSyncer
	if config.Filename != "" {
		// write logs to file
		path := pathlib.NewPath(config.Filename)
		exists, err := path.Exists()
		if err != nil {
			return lm, err
		}
		if !exists {
			err = path.Parent().MkdirAll()
			if err != nil {
				return lm, err
			}
			_, err = path.Create()
			if err != nil {
				return lm, err
			}
		}
		if config.Permissions != 0 {
			err = path.Chmod(fs.FileMode(config.Permissions))
			if err != nil {
				return lm, err
			}
		}

		lm.logFile = &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			LocalTime:  config.LocalTime,
			Compress:   config.Compress,
		}
		writer = zapcore.AddSync(lm.logFile)
	} else {
		// write logs to console
		writer = zapcore.Lock(os.Stdout)
	}

	level := LevelForName(config.Level)
	core := zapcore.NewCore(encoder, writer, level)
	lm.rootLogger = zap.New(
		core,
		// zap.AddStacktrace(zapcore.DebugLevel),

	)

	return lm, nil
}

// NewLogger creates a new child logger.
func (lm LoggerManager) NewLogger(name string, args ...interface{}) *zap.SugaredLogger {
	return lm.rootLogger.Sugar().With(args...).Named(name)
}

// RotateLogFile rotates the underlying log files, if log files are used.
func (lm LoggerManager) RotateLogFile() error {
	if lm.logFile != nil {
		err := lm.logFile.Rotate()
		if err != nil {
			return err
		}
	}
	return nil
}

// Close flushes the logs and closes the underlying log file.
func (lm LoggerManager) Close() error {
	if lm.rootLogger != nil {
		lm.rootLogger.Sync()
		lm.rootLogger = nil
		if lm.logFile != nil {
			logFile := lm.logFile
			lm.logFile = nil
			return logFile.Close()
		}
	}
	return nil
}

// -----------------------------------------------------------------------------

var appLoggerManager LoggerManager

// Initialize initializes the App logger manager.
func Initialize(config Config) (err error) {
	appLoggerManager, err = NewLoggerManager(config)
	return
}

// NewLogger creates a new child logger to be used for general purpose logging
// in the application. It panics if the app logger hasn't been initialized
// beforehand.
func NewLogger(name string, args ...interface{}) *zap.SugaredLogger {
	if appLoggerManager.rootLogger == nil {
		panic("logger has to be initialized first")
	}
	return appLoggerManager.NewLogger(name, args...)
}

// RotateLogFile rotates the underlying App log files, if log files are used.
func RotateLogFile() error {
	return appLoggerManager.RotateLogFile()
}

// Close flushes the App logs and closes the underlying log file.
func Close() error {
	return appLoggerManager.Close()
}

// LevelForName returns the zap log level for the given name.
func LevelForName(name string) zapcore.Level {
	name = strings.ToLower(name)
	switch name {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warning":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "critical":
		return zap.FatalLevel
	}
	return zap.DebugLevel
}
