package command

import (
	"fmt"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	cliML "github.com/jkandasa/jenkinsctl/pkg/model/cli"
	"github.com/jkandasa/jenkinsctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	resourceFile string
)

func init() {
	rootCmd.AddCommand(createResource)
	createResource.PersistentFlags().StringVarP(&resourceFile, "file", "f", "", "resource file")
	err := createResource.MarkPersistentFlagRequired("file")
	if err != nil {
		fmt.Fprintln(ioStreams.ErrOut, "error on fixing a flag", err)
	}
}

var createResource = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from a file",
	Example: `  # create a build using the date in yaml file
  jenkinsctl create -f my_build.yaml
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if CONFIG.JobContext == "" {
			fmt.Fprintf(ioStreams.ErrOut, "job context not set")
			return
		}

		resourceInterface, err := utils.GetResource(resourceFile)
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, err)
			return
		}

		client := jenkins.NewClient(CONFIG, &ioStreams)
		if client == nil {
			return
		}

		switch resource := resourceInterface.(type) {
		case *cliML.KindBuild:
			buildQueueId, err := client.Build(resource.Spec.JobName, resource.Spec.Parameters)
			if err != nil {
				fmt.Fprintln(ioStreams.ErrOut, err)
				return
			}
			fmt.Fprintf(ioStreams.Out, "build created on the job '%s', build queue id:%d\n", resource.Spec.JobName, buildQueueId)
			return

		case *cliML.KindJob:
			jobName, err := client.CreateJob(resource.Spec.JobName, resource.Spec.XMLData)
			if err != nil {
				fmt.Fprintln(ioStreams.ErrOut, err)
				return
			}
			fmt.Fprintf(ioStreams.Out, "job created, job name:%s\n", jobName)
			return

		default:
			fmt.Fprintf(ioStreams.ErrOut, "unknown interface:%T\n", resourceInterface)
			return
		}
	},
}
