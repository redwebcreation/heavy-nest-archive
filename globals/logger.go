package globals

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	loggerConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.Level(Config.Proxy.Logs.Level)),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:      Config.Proxy.Logs.Redirections,
		ErrorOutputPaths: Config.Proxy.Logs.Redirections,
	}

	logger, err := loggerConfig.Build()

	if err != nil {
		panic(err)
	}

	Logger = logger
}
