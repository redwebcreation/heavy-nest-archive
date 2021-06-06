package main

import (
	"fmt"
	"github.com/redwebcreation/hez/cli"
	"github.com/redwebcreation/hez/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	config, _ := core.FindConfig(core.ConfigFile()).Resolve()

	loggerConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.Level(config.Proxy.Logs.Level)),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths: config.Proxy.Logs.Redirections,
	}

	logger, err := loggerConfig.Build()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)

	cli.Execute()
}
