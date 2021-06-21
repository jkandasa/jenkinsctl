package command

import (
	"fmt"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	"github.com/spf13/cobra"
)

var (
	buildNumber          int
	artifactSaveLocation string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.AddCommand(downloadArtifacts)
	downloadArtifacts.PersistentFlags().StringVar(&artifactSaveLocation, "to-dir", "./", "directory to save artifacts")
	downloadArtifacts.PersistentFlags().IntVar(&buildNumber, "build-number", 0, "build number from the job")
	downloadArtifacts.MarkPersistentFlagRequired("build-number")
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download resources from Jenkins server",
}

var downloadArtifacts = &cobra.Command{
	Use:     "artifact",
	Aliases: []string{"artifacts"},
	Short:   "Download artifact of a build",
	Example: `  jenkinsctl download artifact 2101 --to-dir /tmp/artifacts`,
	Run: func(cmd *cobra.Command, args []string) {
		client := jenkins.NewClient(CONFIG)
		if client == nil {
			return
		}

		savedLocation, err := client.DownloadArtifacts(CONFIG.JobContext, buildNumber, artifactSaveLocation)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on downloading artifacts. error:[%s]\n", err)
			return
		}
		fmt.Fprintf(ioStreams.Out, "artifacts are downloaded on the directory: %s\n", savedLocation)
	},
}
