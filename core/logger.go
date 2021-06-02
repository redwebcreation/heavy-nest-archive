package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() *zap.Logger {
	config, _ := GetConfig()

	var outputPaths []string
	var errorOutputPaths []string

	for _, redirection := range config.Logs.Redirections {
		if redirection.For == "out" {
			outputPaths = append(outputPaths, redirection.Value)
		} else {
			errorOutputPaths = append(errorOutputPaths, redirection.Value)
		}
	}

	var loggerConfig zap.Config

	loggerConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(config.Logs.Level))
	loggerConfig.OutputPaths = outputPaths
	loggerConfig.ErrorOutputPaths = errorOutputPaths
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
