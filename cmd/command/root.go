package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	model "github.com/jkandasa/jenkinsctl/pkg/model"
	"github.com/jkandasa/jenkinsctl/pkg/model/config"
	"github.com/jkandasa/jenkinsctl/pkg/printer"
	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ENV_PREFIX       = "JC"
	CONFIG_FILE_NAME = ".jenkinsctl"
	CONFIG_FILE_EXT  = "yaml"
)

var (
	cfgFile   string
	CONFIG    *config.Config  // keep jenkins server details
	ioStreams model.IOStreams // read and write to this stream

	jobContext   string
	hideHeader   bool
	pretty       bool
	outputFormat string

	rootCliLong = `Jenkins Client
  
This client helps you to control your jenkins server from command line.
`
)

var rootCmd = &cobra.Command{
	Use:   "jenkinsctl",
	Short: "Jenkinsctl",
	Long:  rootCliLong,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(ioStreams.Out)
		cmd.SetErr(ioStreams.ErrOut)
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jenkinsctl.yaml)")
	rootCmd.PersistentFlags().StringVarP(&jobContext, "job", "j", "", "Switch to another job")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", printer.OutputConsole, "output format. options: yaml, json, console")
	rootCmd.PersistentFlags().BoolVar(&hideHeader, "hide-header", false, "hides the header on the console output")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty", false, "JSON pretty print")
}

func Execute(streams model.IOStreams) {
	ioStreams = streams
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(ioStreams.ErrOut, err)
		os.Exit(1)
	}
}

func WriteConfigFile() {
	if cfgFile == "" {
		return
	}
	if CONFIG == nil {
		CONFIG = &config.Config{}
	}
	// encode password field
	CONFIG.EncodePassword()

	configBytes, err := yaml.Marshal(CONFIG)
	if err != nil {
		fmt.Fprintf(ioStreams.ErrOut, "error on config file marshal. error:[%s]\n", err.Error())
	}
	err = ioutil.WriteFile(cfgFile, configBytes, os.ModePerm)
	if err != nil {
		fmt.Fprintf(ioStreams.ErrOut, "error on writing config file to disk, filename:%s, error:[%s]\n", cfgFile, err.Error())
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.initConfig
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".jenkinsctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(CONFIG_FILE_NAME)
		viper.SetConfigType(CONFIG_FILE_EXT)

		cfgFile = filepath.Join(home, fmt.Sprintf("%s.%s", CONFIG_FILE_NAME, CONFIG_FILE_EXT))

	}

	viper.SetEnvPrefix(ENV_PREFIX)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
		err = viper.Unmarshal(&CONFIG)
		if err != nil {
			fmt.Fprint(ioStreams.ErrOut, "error on unmarshal of config\n", err)
		}
	} else {
		fmt.Fprint(ioStreams.ErrOut, "error on loading config\n", err)
	}

	if jobContext != "" {
		CONFIG.JobContext = jobContext
	}
}
