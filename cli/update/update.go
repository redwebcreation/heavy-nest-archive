package update

import (
	"encoding/json"
	"fmt"
	box "github.com/redwebcreation/hez/core/embed"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

func run(_ *cobra.Command, _ []string) {
	currentVersion := strings.TrimSpace(string(box.Get("/version")))

	response, err := http.Get("git")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()
	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(data, currentVersion)
}

func NewCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates Hez to the latest version.",
		Run:   run,
	}

	return applyCmd
}
