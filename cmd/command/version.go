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
		client := jenkins.NewClient(CONFIG)
		if client != nil {
			fmt.Fprintln(ioStreams.Out, "Server Version:", client.Version())
		}
	},
}
