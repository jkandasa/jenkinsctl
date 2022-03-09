package command

import (
	"fmt"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	"github.com/jkandasa/jenkinsctl/pkg/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the client and server version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(ioStreams.Out, "Client Version:", version.Get().String())
		if CONFIG.URL == "" {
			fmt.Fprintln(ioStreams.Out, "Server Version: not logged in")
			return
		}
		client, err := jenkins.NewClient(CONFIG, &ioStreams)
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, "error on login", err)
			return
		}
		if client != nil {
			fmt.Fprintln(ioStreams.Out, "Server Version:", client.Version())
		}
	},
}
