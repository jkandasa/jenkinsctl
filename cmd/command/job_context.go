package command

import (
	"fmt"
	"strings"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	"github.com/jkandasa/jenkinsctl/pkg/printer"
	"github.com/spf13/cobra"
)

var (
	depth int
)

func init() {
	rootCmd.AddCommand(jobContextCmd)
	rootCmd.AddCommand(getJobs)

	getJobs.PersistentFlags().IntVar(&depth, "depth", 1, "scan depth of a folder job")
}

var jobContextCmd = &cobra.Command{
	Use:   "job",
	Short: "Switch or set a job",
	Example: `  # set a job
  jenkinsctl set my-another-job
	
  # get the current job
  jenkinsctl job`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintf(ioStreams.ErrOut, "Current job '%s' at '%s'\n", CONFIG.JobContext, CONFIG.URL)
			return
		}
		client, err := jenkins.NewClient(CONFIG, &ioStreams)
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, "error on login", err)
			return
		}
		if client != nil {
			CONFIG.JobContext = strings.TrimSpace(args[0])
			WriteConfigFile()
			fmt.Fprintf(ioStreams.ErrOut, "Switched to '%s' at '%s'\n", CONFIG.JobContext, CONFIG.URL)
		}
	},
}

var getJobs = &cobra.Command{
	Use:   "jobs",
	Short: "Display existing jobs",
	Example: `  # display existing jobs
  jenkinsctl jobs

	# display existing jobs with depth flag
  jenkinsctl jobs --depth 2`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := jenkins.NewClient(CONFIG, &ioStreams)
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, "error on login", err)
			return
		}
		if client == nil {
			return
		}
		jobs, err := client.ListJobs(depth)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on listing jobs. error:[%s]", err.Error())
		}

		headers := []string{"color", "name", "class", "url"}
		data := make([]interface{}, 0)
		for _, job := range jobs {
			data = append(data, Job{Name: job.Name, Class: job.Class, URL: job.Url, Color: job.Color})
		}
		printer.Print(ioStreams.Out, headers, data, hideHeader, outputFormat, pretty)
	},
}
