package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NetworkCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "network",
		Short:                 "Manage networks",
		Long:                  `Manage networks`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The network main command have been executed")
		},
	}

	cmd.AddCommand(NetworkCreateCommand())
	cmd.AddCommand(NetworkConnectCommand())
	cmd.AddCommand(NetworkDisconnectCommand())
	cmd.AddCommand(NetworkListCommand())
	cmd.AddCommand(NetworkRemoveCommand())
	return cmd
}

func NetworkCreateCommand() *cobra.Command {}
