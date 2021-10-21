package node

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runListenCommand(_ *cobra.Command, _ []string) error {
	r := gin.Default()

	r.GET("/join", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I'm a node!",
		})
	})

	r.Run()
}

func ListenCommand() *cobra.Command {
	return cmd.CreateCommand(&cobra.Command{
		Use:   "listen",
		Short: "Starts the node server",
	}, nil, runListenCommand)
}
