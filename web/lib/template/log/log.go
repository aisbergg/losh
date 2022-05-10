package log

import (
	"losh/lib/log"

	"github.com/rotisserie/eris"
	"go.uber.org/zap"
)

var appLogger log.LoggerWrapper
var accessLogger log.LoggerWrapper

// Initialize initializes the application and access logging.
func Initialize(appLogCfg AppLogConfig, accLogCfg AccessLogConfig) error {
	var err error

	// init application logger
	appLogger, err = log.CreateLogger(log.LevelForName(appLogCfg.Level), log.CommonLogConfig{
		Format: appLogCfg.Format,
		FileLogConfig: log.FileLogConfig{
			Filename:    appLogCfg.Filename,
			Rotate:      appLogCfg.Rotate,
			MaxSize:     appLogCfg.MaxSize,
			MaxAge:      appLogCfg.MaxAge,
			MaxBackups:  appLogCfg.MaxBackups,
			LocalTime:   appLogCfg.LocalTime,
			Compress:    appLogCfg.Compress,
			Permissions: appLogCfg.Permissions,
		},
	})
	if err != nil {
		return eris.Wrap(err, "failed to initialize application logging")
	}

	// init access logger
	logLevel := zap.FatalLevel
	if accLogCfg.Enabled {
		logLevel = zap.InfoLevel
	}
	accessLogger, err = log.CreateLogger(logLevel, log.CommonLogConfig{
		Format: accLogCfg.Format,
		FileLogConfig: log.FileLogConfig{
			Filename:    accLogCfg.Filename,
			Rotate:      accLogCfg.Rotate,
			MaxSize:     accLogCfg.MaxSize,
			MaxAge:      accLogCfg.MaxAge,
			MaxBackups:  accLogCfg.MaxBackups,
			LocalTime:   accLogCfg.LocalTime,
			Compress:    accLogCfg.Compress,
			Permissions: accLogCfg.Permissions,
		},
	})
	if err != nil {
		return eris.Wrap(err, "failed to initialize access logging")
	}

	return nil
}

// NewLogger creates a new child logger to be used for general purpose logging
// in the application. It panics if the app logger hasn't been initialized
// beforehand.
func NewLogger(name string, args ...interface{}) *zap.SugaredLogger {
	if appLogger.Logger == nil {
		panic("logger has to be initialized first")
	}
	return appLogger.Logger.Sugar().With(args...).Named(name)
}

// AccessLogger returns the logger for logging web access requests. It panics if
// the app logger hasn't been initialized beforehand.
func AccessLogger() *zap.Logger {
	if accessLogger.Logger == nil {
		panic("logger has to be initialized first")
	}
	return accessLogger.Logger
}

// RotateLogFile rotates the underlying log files.
func RotateLogFile() error {
	if appLogger.LogFile != nil {
		err := appLogger.LogFile.Rotate()
		if err != nil {
			return err
		}
	}
	if accessLogger.LogFile != nil {
		err := accessLogger.LogFile.Rotate()
		if err != nil {
			return err
		}
	}
	return nil
}
