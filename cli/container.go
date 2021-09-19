package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

type containerCreateOptions struct {
	name       string
	network    string
	mountDevfs bool
	env        []string
	volume     []string
	jailParam  []string
}

func newContainerCreateCommand() *cobra.Command {
	opts := containerCreateOptions{}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new container",
		Long:  `Create a new container loooong`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container create command has been executed")
			fmt.Println(args)
			fmt.Println(opts)
		},
	}

	flags := createCmd.Flags()
	flags.BoolVar(&opts.mountDevfs, "mount.devfs", true, "Toggle devfs mount")
	flags.StringVar(&opts.name, "name", "", "Assign a name to the container")
	flags.StringVar(&opts.network, "network", "", "Connect a container to a network")
	flags.StringSliceVarP(&opts.volume, "volume", "v", []string{""}, "Bind mount a volume to the container")
	flags.StringSliceVarP(&opts.env, "env", "e", []string{""}, "Set environment variables (e.g. --env FIRST=env --env SECOND=env)")
	flags.StringSliceVarP(&opts.jailParam, "jailparam", "J", []string{""}, "Specify a jail parameter (see jail(8) for details)")
	return createCmd
}

func newContainerRemoveCommand() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove one or more containers",
		Long:  `Remove one or more containers loooong`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container rm command has been executed")
		},
	}
	return removeCmd
}

func newContainerStartCommand() *cobra.Command {
	var attach bool
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start one or more stopped containers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container rm command has been executed")
		},
	}

	cmd.Flags().BoolVarP(&attach, "attach", "a", true, "Attach STDOUT/STDERR")
	return cmd
}

func newContainerStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop one or more running containers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container stop command has been executed")
		},
	}
	return cmd
}

func newContainerListCommand() *cobra.Command {
	var all bool

	listCmd := &cobra.Command{
		Use:   "ls",
		Short: "List containers",
		Long:  `List containers loooong`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container ls command has been executed")
		},
	}
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (default shows just running)")
	return listCmd
}

func newContainerCommand() *cobra.Command {
	containerCmd := &cobra.Command{
		Use:                   "container",
		Short:                 "Manage containers",
		Long:                  `Manage containers loooong`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container main command have been executed")
		},
	}

	containerCmd.AddCommand(newContainerCreateCommand())
	containerCmd.AddCommand(newContainerRemoveCommand())
	containerCmd.AddCommand(newContainerStartCommand())
	containerCmd.AddCommand(newContainerStopCommand())
	containerCmd.AddCommand(newContainerListCommand())
	return containerCmd
}
