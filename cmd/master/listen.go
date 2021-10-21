package master

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wormable/nest/cmd"
)

func runListenCommand(_ *cobra.Command, _ []string) error {
	r := gin.Default()

	r.GET("/join", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "I'm the master!",
		})
	})

	r.Run()
}

func ListenCommand() *cobra.Command {
	return cmd.CreateCommand(&cobra.Command{
		Use:   "listen",
		Short: "Start the master server",
	}, nil, runListenCommand)
}
