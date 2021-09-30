//HOWTOS:
//$ go get -u github.com/spf13/cobra

package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	version bool
	debug   bool
	host    string

	RootCmd = &cobra.Command{
		Use:     "jcli",
		Short:   "A cli-tool for jocker",
		Long:    `JCli is the reference cli-tool for interacting with jocker-engine`,
		Version: "0.0.1",
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	RootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "Daemon socket to connect to: tcp://[host]:[port][path] or unix://[/path/to/socket]")
	RootCmd.AddCommand(NewContainerCommand())
}
