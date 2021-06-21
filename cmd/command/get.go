package command

import (
	"fmt"
	"strconv"
	"time"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	"github.com/jkandasa/jenkinsctl/pkg/printer"
	"github.com/spf13/cobra"
)

type Job struct {
	Name  string `json:"name" yaml:"name" structs:"name"`
	Class string `json:"class" yaml:"class" structs:"class"`
}

// JobParameters struct
type JobParameters struct {
	Name         string      `json:"name" yaml:"name" structs:"name"`
	Type         string      `json:"type" yaml:"type" structs:"type"`
	DefaultName  string      `json:"defaultName" yaml:"default_name" structs:"default name"`
	DefaultValue interface{} `json:"defaultValue" yaml:"default_value" structs:"default value"`
	Description  string      `json:"description" yaml:"description" structs:"description"`
}

type TableData struct {
	Key   string `json:"key" yaml:"key" structs:"key"`
	Value string `json:"value" yaml:"value" structs:"value"`
}

var (
	limit int
	watch bool
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getBuilds)
	getCmd.AddCommand(getParameters)
	getCmd.AddCommand(getConsole)

	getCmd.PersistentFlags().IntVar(&limit, "limit", 5, "limit the number of entries to display")
	getConsole.PersistentFlags().BoolVarP(&watch, "watch", "w", false, "watch build console logs")

}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Example: `  # get builds
  jenkinsctl get builds

  # get parameters
  jenkinsctl get parameters`,
}

var getBuilds = &cobra.Command{
	Use:     "build",
	Aliases: []string{"builds"},
	Short:   "Displays builds of a job",
	Example: `  # get builds
  jenkinsctl get builds

  # get builds on a different job (temporary switch)
  jenkinsctl get builds -j my-another-job

  # get limited builds
  jenkinsctl get builds --limit 2

  # get output as yaml
  jenkinsctl get builds --limit 2 --output yaml

  # get output as json
  jenkinsctl get builds --limit 2 --output json --pretty

  # get a particular build details with build number
  jenkinsctl get build 61`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := jenkins.NewClient(CONFIG)
		if client == nil {
			return
		}

		if len(args) > 0 {
			buildNumber, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(ioStreams.ErrOut, "build should be an integer number", err)
				return
			}
			build, err := client.GetBuild(CONFIG.JobContext, buildNumber, false)
			if err != nil {
				fmt.Fprintln(ioStreams.ErrOut, "error on getting build details", err)
				return
			}
			if outputFormat != printer.OutputConsole {
				printer.Print(ioStreams.Out, nil, build, false, outputFormat, pretty)
				return
			}

			headers := []string{"key", "value"}
			rows := make([]interface{}, 0)
			rows = append(rows, TableData{Key: "URL", Value: build.URL})
			rows = append(rows, TableData{Key: "Build Number", Value: fmt.Sprintf("%d", build.Number)})
			rows = append(rows, TableData{Key: "Triggered By", Value: build.TriggeredBy})
			rows = append(rows, TableData{Key: "Result", Value: build.Result})
			rows = append(rows, TableData{Key: "Is Running", Value: fmt.Sprintf("%v", build.IsRunning)})
			rows = append(rows, TableData{Key: "Duration", Value: (time.Duration(build.TestResult.Duration) * time.Second).String()})
			rows = append(rows, TableData{Key: "Revision", Value: build.Revision})
			rows = append(rows, TableData{Key: "Revision Branch", Value: build.RevisionBranch})
			rows = append(rows, TableData{Key: "Timestamp", Value: build.Timestamp.String()})
			printer.Print(ioStreams.Out, headers, rows, false, outputFormat, pretty)

			fmt.Fprintf(ioStreams.Out, "\nParameters:\n")
			headers = []string{"name", "value"}
			rows = make([]interface{}, 0)
			for _, parameter := range build.Parameters {
				rows = append(rows, parameter)
			}
			printer.Print(ioStreams.Out, headers, rows, true, outputFormat, pretty)

			fmt.Fprintf(ioStreams.Out, "\nTest Result:\n")
			if !build.TestResult.Empty {
				headers = []string{"key", "value"}
				rows = make([]interface{}, 0)
				rows = append(rows, TableData{Key: "Passed", Value: fmt.Sprintf("%d", build.TestResult.PassCount)})
				rows = append(rows, TableData{Key: "Failed", Value: fmt.Sprintf("%d", build.TestResult.FailCount)})
				rows = append(rows, TableData{Key: "Skipped", Value: fmt.Sprintf("%d", build.TestResult.SkipCount)})
				rows = append(rows, TableData{Key: "Duration", Value: (time.Duration(build.TestResult.Duration) * time.Second).String()})
				printer.Print(ioStreams.Out, headers, rows, false, outputFormat, pretty)
			}

			fmt.Fprintf(ioStreams.Out, "\nArtifacts:\n")
			headers = []string{"path"}
			rows = make([]interface{}, 0)
			for _, artifact := range build.Artifacts {
				rows = append(rows, artifact)
			}
			printer.Print(ioStreams.Out, headers, rows, false, outputFormat, pretty)

		} else {
			builds, err := client.ListBuilds(CONFIG.JobContext, limit, false)
			if err != nil {
				fmt.Fprintf(ioStreams.ErrOut, "error on listing builds. error:[%s]", err.Error())
			}

			headers := []string{"number", "triggered by", "result", "is running", "duration", "timestamp", "revision"}
			data := make([]interface{}, 0)
			for _, build := range builds {
				data = append(data, build)
			}
			printer.Print(ioStreams.Out, headers, data, hideHeader, outputFormat, pretty)
		}
	},
}

var getParameters = &cobra.Command{
	Use:     "parameter",
	Aliases: []string{"parameters"},
	Short:   "Displays all the parameters of a job",
	Example: `  # get parameters
  jenkinsctl get parameters

  # get parametes as yaml
  jenkinsctl get parameters --output yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		client := jenkins.NewClient(CONFIG)
		if client == nil {
			return
		}

		parameters, err := client.ListParameters(CONFIG.JobContext)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on listing parameters. error:[%s]", err.Error())
		}
		headers := []string{"name", "default value", "type", "description"}
		data := make([]interface{}, 0)
		for _, param := range parameters {
			data = append(data,
				JobParameters{
					Name:         param.Name,
					Type:         param.Type,
					DefaultName:  param.DefaultParameterValue.Name,
					DefaultValue: param.DefaultParameterValue.Value,
					Description:  param.Description,
				})
		}
		printer.Print(ioStreams.Out, headers, data, hideHeader, outputFormat, pretty)
	},
}

var getConsole = &cobra.Command{
	Use:   "console",
	Short: "Print the console logs for a build in a job",
	Example: `  # get console output of a build
  jenkinsctl get console 61

  # watch a console output of a build
  jenkinsctl get console 61 --watch`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		buildNumber, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintln(ioStreams.ErrOut, "expecting a build number, error:", err)
			return
		}

		// old
		client := jenkins.NewClient(CONFIG)
		if client == nil {
			return
		}

		consoleLog, err := client.GetConsole(CONFIG.JobContext, buildNumber, watch, ioStreams.Out)
		if err != nil {
			fmt.Fprintf(ioStreams.ErrOut, "error on listing build console. error:[%s]", err.Error())
			return
		}

		fmt.Fprint(ioStreams.Out, consoleLog)
	},
}
