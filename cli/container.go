package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	Openapi "jcli/client"

	"github.com/spf13/cobra"
)

const url = "http://localhost:8085/"

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

func NewContainerCreateCommand() *cobra.Command {
	config := Openapi.ContainerCreateJSONRequestBody{
		Networks:  &([]string{}),
		Volumes:   &([]string{}),
		Env:       &([]string{}),
		JailParam: &([]string{}),
	}

	var name string

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new container",
		Long:  `Create a new container loooong`,
		Run: func(cmd *cobra.Command, args []string) {
			RunContainerCreate(cmd, &name, config, args)
		},
	}

	flags := createCmd.Flags()
	flags.StringVar(&name, "name", "", "Assign a name to the container")
	flags.StringSliceVar(config.Networks, "network", []string{}, "Connect a container to a network")
	flags.StringSliceVarP(config.Volumes, "volume", "v", []string{}, "Bind mount a volume to the container")
	flags.StringSliceVarP(config.Env, "env", "e", []string{}, "Set environment variables (e.g. --env FIRST=env --env SECOND=env)")
	flags.StringSliceVarP(config.JailParam, "jailparam", "J", []string{"mount.devfs"}, "Specify a jail parameter (see jail(8) for details)")
	return createCmd
}

func RunContainerCreate(cmd *cobra.Command, name *string, body Openapi.ContainerCreateJSONRequestBody, args []string) {
	container_cmd := args[1:]
	image := args[0]
	body.Cmd = &container_cmd
	body.Image = &image

	params := Openapi.ContainerCreateParams{}
	if *name != "" {
		params = Openapi.ContainerCreateParams{Name: name}
	}

	client := NewHTTPClient()

	response, _ := client.ContainerCreateWithResponse(context.TODO(), &params, body)
	if response.StatusCode() != 201 {
		fmt.Println("Jocker engine returned unsuccesful statuscode: ", response.Status())
		os.Exit(0)
	}
	fmt.Println(response.JSON201.Id)
}

func NewContainerRemoveCommand() *cobra.Command {
	removeCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove one or more containers",
		Long:  `Remove one or more containers loooong`,
		Run: func(cmd *cobra.Command, args []string) {
			RunContainerRemove(cmd, args)
		},
	}
	return removeCmd
}

func RunContainerRemove(cmd *cobra.Command, args []string) {
	container_id := args[0]
	client := NewHTTPClient()
	response, _ := client.ContainerDeleteWithResponse(context.TODO(), container_id)
	status_code := response.StatusCode()

	switch {
	case status_code == 204:
		fmt.Println("succesfully removed container")
	case status_code == 404:
		fmt.Println("no such container")
	case status_code == 500:
		fmt.Println("internal server error")
	}
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
			RunContainerList(cmd, all, args)
		},
	}
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (default shows just running)")
	return listCmd
}

func RunContainerList(cmd *cobra.Command, all bool, args []string) {
	params := Openapi.ContainerListParams{
		All: &all,
	}

	client := NewHTTPClient()

	response, err := client.ContainerListWithResponse(context.TODO(), &params)
	if err != nil {
		fmt.Println("Could not connect to jocker engine daemon: ", err)
		return
	}

	if response.StatusCode() != 200 {
		fmt.Println("Jocker engine returned non-200 statuscode: ", response.Status())
		return
	}
	PrintContainerList(response.JSON200)
}

func PrintContainerList(container_list *[]Openapi.ContainerSummary) {
	fmt.Println(
		Cell("CONTAINER ID", 12), Sp(3),
		Cell("IMAGE", 15), Sp(3),
		Cell("COMMAND", 23), Sp(3),
		Cell("CREATED", 18), Sp(3),
		Cell("STATUS", 7), Sp(3),
		"NAME",
	)

	var running string

	for _, c := range *container_list {
		if *c.Running {
			running = "running"
		} else {
			running = "stopped"
		}
		created, _ := time.Parse(time.RFC3339, *c.Created)
		since_created := time.Since(created)

		fmt.Println(
			Cell(*c.Id, 12), Sp(1),
			Cell(*c.ImageId, 15), Sp(1),
			Cell(*c.Command, 23), Sp(1),
			Cell(HumanDuration(since_created)+" ago", 18), Sp(1),
			Cell(running, 7), Sp(1),
			*c.Name,
		)
	}
}

func NewHTTPClient() *Openapi.ClientWithResponses {
	client, err := Openapi.NewClientWithResponses(url)
	if err != nil {
		fmt.Println("Internal error: ", err)
		os.Exit(1)
	}
	return client
}

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.).
func HumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds == 1 {
		return "1 second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours() + 0.5); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	} else if hours < 24*7*2 {
		return fmt.Sprintf("%d days", hours/24)
	} else if hours < 24*30*2 {
		return fmt.Sprintf("%d weeks", hours/24/7)
	} else if hours < 24*365*2 {
		return fmt.Sprintf("%d months", hours/24/30)
	}
	return fmt.Sprintf("%d years", int(d.Hours())/24/365)
}

func Cell(word string, max_len int) string {
	word_length := len(word)

	if word_length <= max_len {
		return word + Sp(max_len-word_length) + Sp(2)
	} else {
		return word[:max_len] + ".."
	}
}

func Sp(n int) string {
	return strings.Repeat(" ", n)
}
