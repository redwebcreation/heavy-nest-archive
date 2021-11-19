package main

import (
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
	"github.com/wormable/nest/common"
	"github.com/wormable/ui"
	"log/syslog"
	"os"
	"strconv"
)

func main() {
	policy := common.LogPolicy{
		Name: "default",
		Rules: []common.LogRule{
			{
				Level: "debug",
				Redirections: []common.LogRedirection{
					{
						Type:     "syslog",
						Facility: strconv.Itoa(int(syslog.LOG_SYSLOG)),
					},
				},
			},
		},
	}

	policy.Log(syslog.LOG_ERR, "test")

	os.Exit(0)
	cli := &cobra.Command{
		Use:   "nest",
		Short: "nest makes orchestrating containers easy.",
		Long:  "nest is to tool to orchestrate containers and manage the environment around them.",
	}

	cli.AddCommand(
		cmd.ApplyCommand(),
		cmd.DiagnoseCommand(),
		cmd.StopCommand(),
		cmd.SelfUpdateCommand(),
		cmd.ProxyCommand(),
		cmd.InitCommand(),
	)

	err := cli.Execute()
	ui.Check(err)
}
