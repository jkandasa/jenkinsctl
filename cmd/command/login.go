package command

import (
	"fmt"

	jenkins "github.com/jkandasa/jenkinsctl/pkg/jenkins"
	"github.com/spf13/cobra"
)

var (
	username              string
	password              string
	insecureSkipTLSVerify bool
)

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username to login with jenkins server")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password or token to login with jenkins server")
	loginCmd.Flags().BoolVar(&insecureSkipTLSVerify, "insecure-skip-tls-verify", false,
		"If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure")

	rootCmd.AddCommand(logoutCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a server",
	Example: `  # login to the server with username and password/token
  jenkinsctl login http://localhost:8080 --username jenkins --password my_token

  # login to the insecure server (with SSL certificate)
  jenkinsctl login https://localhost:8443 --username jenkins --password my_token  --insecure-skip-tls-verify`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CONFIG.URL = args[0]
		CONFIG.Username = username
		CONFIG.Password = password
		CONFIG.InsecureSkipTLSVerify = insecureSkipTLSVerify
		client := jenkins.NewClient(CONFIG)
		if client != nil {
			fmt.Fprintln(ioStreams.ErrOut, "Login successful.")
			WriteConfigFile()
		}
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from a server",
	Example: `  # logout from a server
  jenkinsctl logout`,
	Run: func(cmd *cobra.Command, args []string) {
		CONFIG.URL = ""
		CONFIG.Username = ""
		CONFIG.Password = ""
		CONFIG.InsecureSkipTLSVerify = false
		CONFIG.JobContext = ""
		fmt.Fprintln(ioStreams.ErrOut, "Logout successful.")
		WriteConfigFile()
	},
}
