package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() *zap.Logger {
	config, _ := FindConfig(ConfigFile()).Resolve()

	var loggerConfig zap.Config

	loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(config.Logs.Level))
	loggerConfig.OutputPaths = config.Logs.Redirections
	loggerConfig.Encoding = "json"
	loggerConfig.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}

	logger, _ := loggerConfig.Build()

	defer logger.Sync()

	return logger
}
