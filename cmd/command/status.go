package command

import (
	"fmt"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	"github.com/jkandasa/jenkinsctl/pkg/printer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays an overview of the jenkins server",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := jenkins.NewClient(CONFIG, &ioStreams)
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, "error on login", err)
			return
		}
		if client == nil {
			return
		}
		response, err := client.Status()
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, err)
			return
		}
		if outputFormat != printer.OutputConsole {
			printer.Print(ioStreams.Out, nil, response, hideHeader, outputFormat, pretty)
			return
		}

		fmt.Fprintf(ioStreams.Out, "Description: %s\n", response.Description)
		fmt.Fprintf(ioStreams.Out, "Number of executers: %d\n", response.NumExecutors)
		fmt.Fprintf(ioStreams.Out, "\nLabels:\n")
		for _, label := range response.AssignedLabels {
			for key, value := range label {
				fmt.Fprintf(ioStreams.Out, "\t%s: %s\n", key, value)
			}
		}

		fmt.Fprintf(ioStreams.Out, "\nNode:\n")
		fmt.Fprintf(ioStreams.Out, "\tName: %s\n", response.NodeName)
		fmt.Fprintf(ioStreams.Out, "\tDescription: %s\n", response.NodeDescription)

		fmt.Fprint(ioStreams.Out, "\nJobs:\n")

		headers := []string{"name", "class"}
		data := make([]interface{}, 0)
		for _, job := range response.Jobs {
			data = append(data, Job{Class: job.Class, Name: job.Name})
		}
		printer.Print(ioStreams.Out, headers, data, hideHeader, outputFormat, pretty)
	},
}
