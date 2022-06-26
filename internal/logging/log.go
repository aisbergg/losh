package logging

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerWrapper is just for wrapping the logging instance together with the log
// file.
type LoggerWrapper struct {
	Logger  *zap.Logger
	LogFile *lumberjack.Logger
}

// CreateLogger creates a new logger with the given configuration.
func CreateLogger(level zapcore.Level, cfg CommonLogConfig) (LoggerWrapper, error) {
	var err error
	var lw LoggerWrapper

	if cfg.Filename == "" {
		// console logger
		logCfg := zap.NewProductionConfig()
		logCfg.Level = zap.NewAtomicLevelAt(level)
		logCfg.Encoding = "console"
		lw.Logger, err = logCfg.Build()
		if err != nil {
			return lw, err
		}

	} else {
		// file logger

		// TODO: file checks
		// TODO: file permission
		lw.LogFile = &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			LocalTime:  cfg.LocalTime,
			Compress:   cfg.Compress,
		}
		logWriter := zapcore.AddSync(lw.LogFile)
		var logEncoder zapcore.Encoder
		if cfg.Format == "json" {
			logEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		} else {
			logEncoder = zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
		}
		lw.Logger = zap.New(zapcore.NewCore(
			logEncoder,
			logWriter,
			level,
		))
	}

	defer lw.Logger.Sync()

	return lw, nil
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
