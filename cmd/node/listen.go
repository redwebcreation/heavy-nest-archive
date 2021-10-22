package node

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
	"github.com/wormable/nest/globals"
)

func runListenCommand(_ *cobra.Command, _ []string) error {
	r := gin.Default()

	r.GET("/version", func(c *gin.Context) {
		c.String(200, "nest@%s", globals.Version)
	})

	return r.Run(":80")
}

func ListenCommand() *cobra.Command {
	return cmd.CreateCommand(&cobra.Command{
		Use:   "listen",
		Short: "Starts the node server",
	}, nil, runListenCommand)
}
