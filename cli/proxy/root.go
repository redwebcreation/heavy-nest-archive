package proxy

import (
	"fmt"
	"github.com/redwebcreation/hez/core"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func NewCommand() *cobra.Command {
	proxyCommand := &cobra.Command{
		Use:   "proxy",
		Short: "Manage the reverse proxy",
	}

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

	proxyCommand.AddCommand(initEnableCommand())
	proxyCommand.AddCommand(initDisableCommand())
	proxyCommand.AddCommand(initStatusCommand())
	proxyCommand.AddCommand(initRunCommand())

	return proxyCommand
}
