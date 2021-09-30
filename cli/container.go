package cli

import (
	"context"
	"fmt"
	"io"

	Openapi "jcli/client"

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

func NewContainerCreateCommand() *cobra.Command {
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

func NewContainerRemoveCommand() *cobra.Command {
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

func NewContainerStartCommand() *cobra.Command {
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

func NewContainerStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop one or more running containers",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container stop command has been executed")
		},
	}
	return cmd
}

func NewContainerListCommand() *cobra.Command {
	var all bool

	listCmd := &cobra.Command{
		Use:   "ls",
		Short: "List containers",
		Long:  `List containers loooong`,
		Run: func(cmd *cobra.Command, args []string) {
			RunContainerList(cmd, args)
		},
	}
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (default shows just running)")
	return listCmd
}

func RunContainerList(cmd *cobra.Command, args []string) {
	// FIXME: Simple and messy PoC (that works!)
	client, err := Openapi.NewClient("http://localhost:8085/")
	fmt.Println("What", err)
	if err != nil {
		fmt.Println("Internal error: ", err)
	}
	all := true
	params := Openapi.ContainerListParams{
		All: &all,
	}
	response, err := client.ContainerList(context.TODO(), &params)
	if err != nil {
		fmt.Println("Could not connect to jocker engine daemon: ", err)
	}
	fmt.Println("The container ls command has been executed", response)
	bytes, err := io.ReadAll(response.Body)
	fmt.Println("Body content:", string(bytes))
}

func NewContainerCommand() *cobra.Command {
	containerCmd := &cobra.Command{
		Use:                   "container",
		Short:                 "Manage containers",
		Long:                  `Manage containers loooong`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container main command have been executed")
		},
	}

	containerCmd.AddCommand(NewContainerCreateCommand())
	containerCmd.AddCommand(NewContainerRemoveCommand())
	containerCmd.AddCommand(NewContainerStartCommand())
	containerCmd.AddCommand(NewContainerStopCommand())
	containerCmd.AddCommand(NewContainerListCommand())
	return containerCmd
}
