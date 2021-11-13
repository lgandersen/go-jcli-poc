package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	Openapi "jcli/client"

	"github.com/spf13/cobra"
)

const jocker_engine_url = "http://localhost:8085/"
const ws_container_attach = "ws://localhost:8085/containers/%s/attach"
const succesful_ws_exit = "websocket: close 1000 (normal): exit:"

func ContainerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "container",
		Short:                 "Manage containers",
		Long:                  `Manage containers`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("The container main command have been executed")
		},
	}

	cmd.AddCommand(ContainerCreateCommand())
	cmd.AddCommand(ContainerRemoveCommand())
	cmd.AddCommand(ContainerStartCommand())
	cmd.AddCommand(ContainerStopCommand())
	cmd.AddCommand(ContainerListCommand())
	return cmd
}

func ContainerCreateCommand() *cobra.Command {
	config := Openapi.ContainerCreateJSONRequestBody{
		Networks:  &([]string{}),
		Volumes:   &([]string{}),
		Env:       &([]string{}),
		JailParam: &([]string{}),
	}

	var name string

	cmd := &cobra.Command{
		Use:                   "create [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short:                 "Create a new container",
		Long:                  `Create a new container`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := PostContainerCreate(&name, config, args)
			if err == nil {
				fmt.Println(response.JSON201.Id)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&name, "name", "", "Assign a name to the container")
	flags.StringSliceVar(config.Networks, "network", []string{}, "Connect a container to a network")
	flags.StringSliceVarP(config.Volumes, "volume", "v", []string{}, "Bind mount a volume to the container")
	flags.StringSliceVarP(config.Env, "env", "e", []string{}, "Set environment variables (e.g. --env FIRST=env --env SECOND=env)")
	flags.StringSliceVarP(config.JailParam, "jailparam", "J", []string{"mount.devfs"}, "Specify a jail parameter, see jail(8) for details")
	return cmd
}

func PostContainerCreate(name *string, body Openapi.ContainerCreateJSONRequestBody, args []string) (*Openapi.ContainerCreateResponse, error) {
	container_cmd := args[1:]
	image := args[0]
	body.Cmd = &container_cmd
	body.Image = &image

	params := Openapi.ContainerCreateParams{}
	if *name != "" {
		params = Openapi.ContainerCreateParams{Name: name}
	}

	client := NewHTTPClient()

	response, err := client.ContainerCreateWithResponse(context.TODO(), &params, body)
	if err != nil {
		fmt.Println("Could not connect to jocker engine daemon: ", err)
		return response, err
	}

	if response.StatusCode() != 201 {
		fmt.Println("Jocker engine returned unsuccesful statuscode: ", response.Status())
		return response, errors.New("non-200 statuscode")
	}
	return response, nil
}

func ContainerRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rm",
		Short:                 "Remove one or more containers",
		Long:                  `Remove one or more containers`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := PostContainerRemove(args)
			if err == nil {
				fmt.Println(response.JSON200.Id)
			}
		},
	}
	return cmd
}

func PostContainerRemove(args []string) (*Openapi.ContainerDeleteResponse, error) {
	container_id := args[0]
	client := NewHTTPClient()
	response, _ := client.ContainerDeleteWithResponse(context.TODO(), container_id)
	status_code := response.StatusCode()

	switch {
	case status_code == 200:
		//fmt.Println("succesfully removed container")
		return response, nil
	case status_code == 404:
		return response, errors.New("no such container")
	case status_code == 500:
		return response, errors.New("internal server error")
	default:
		return response, errors.New("unknown status-code received from jocker engine: " + response.Status())
	}
}

func ContainerStartCommand() *cobra.Command {
	var attach bool
	cmd := &cobra.Command{
		Use:                   "start [OPTIONS] CONTAINER [CONTAINER...]",
		Short:                 "Start one or more stopped containers",
		Long:                  "Start one or more stopped containers. Attach to STDOUT/STDERR if only one container is started",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			if attach {
				StartAndAttachToContainer(args)
			} else {
				StartSeveralContainers(args)
			}
		},
	}

	cmd.Flags().BoolVarP(&attach, "attach", "a", true, "Attach STDOUT/STDERR")
	return cmd
}

func StartSeveralContainers(args []string) []string {
	client := NewHTTPClient()
	container_ids := make([]string, len(args))
	var container_id string
	for i, container := range args {
		container_id = StartSingleContainer(client, container)
		container_ids[i] = container_id
		if container_id != "" {
			fmt.Println(container_id)
		}
	}
	return container_ids
}

func StartAndAttachToContainer(args []string) {
	if len(args) != 1 {
		fmt.Println("When attaching to STDOUT/STDERR only 1 container can be started")
		return
	}
	container_id := args[0]
	endpoint := fmt.Sprintf(ws_container_attach, container_id)

	done, interrupt, ws := Dial(endpoint)
	go ListenForWSMessages(done, ws)
	StartSingleContainer(NewHTTPClient(), container_id)
	AwaitDoneOrUserInterrupt(done, interrupt, ws)

}

func StartSingleContainer(client *Openapi.ClientWithResponses, container string) string {
	response, err := client.ContainerStartWithResponse(context.TODO(), container)
	status_code := response.StatusCode()
	switch {
	case err != nil:
		fmt.Println("error requesting container start: ", err)

	case status_code == 200 && response.JSON200 == nil:
		fmt.Println("could not parse jocker engine response")

	case status_code == 200:
		return response.JSON200.Id

	case status_code == 304:
		fmt.Println("container already started")

	case status_code == 404:
		fmt.Println("no such container")

	case status_code == 500:
		fmt.Println("internal server error")

	default:
		fmt.Println("unknown status-code received from jocker engine: ", response.Status())
	}
	return ""
}

func ContainerStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "stop",
		Short:                 "Stop one or more running containers",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := ContainerStop(args)
			if err != nil {
				fmt.Println(response.JSON200.Id)
			}
		},
	}
	return cmd
}

func ContainerStop(args []string) (*Openapi.ContainerStopResponse, error) {
	name_or_id := args[0]
	client := NewHTTPClient()
	response, err := client.ContainerStopWithResponse(context.TODO(), name_or_id)
	if err != nil {
		return response, err
	}

	status_code := response.StatusCode()

	switch {
	case status_code == 200:
		return response, nil

	case status_code == 304:
		return response, errors.New("container already stopped")

	case status_code == 404:
		return response, errors.New("no such container")

	case status_code == 500:
		return response, errors.New("internal server error")

	default:
		return response, errors.New("unknown status-code received from jocker engine: " + string(response.Status()))
	}
}

func ContainerListCommand() *cobra.Command {
	var all bool

	listCmd := &cobra.Command{
		Use:                   "ls",
		Short:                 "List containers",
		Long:                  `List containers`,
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := GetContainerList(all)
			if err == nil {
				PrintContainerList(response.JSON200)
			}
		},
	}
	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (default shows just running)")
	return listCmd
}

func GetContainerList(all bool) (*Openapi.ContainerListResponse, error) {
	params := Openapi.ContainerListParams{
		All: &all,
	}

	client := NewHTTPClient()

	response, err := client.ContainerListWithResponse(context.TODO(), &params)
	if err != nil {
		fmt.Println("Could not connect to jocker engine daemon: ", err)
		return response, err
	}

	if response.StatusCode() != 200 {
		fmt.Println("Jocker engine returned non-200 statuscode: ", response.Status())
		return response, errors.New("non-200 statuscode")
	}
	return response, nil
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

	for _, ws := range *container_list {
		if *ws.Running {
			running = "running"
		} else {
			running = "stopped"
		}
		created, _ := time.Parse(time.RFC3339, *ws.Created)
		since_created := time.Since(created)

		fmt.Println(
			Cell(*ws.Id, 12), Sp(1),
			Cell(*ws.ImageId, 15), Sp(1),
			Cell(*ws.Command, 23), Sp(1),
			Cell(HumanDuration(since_created)+" ago", 18), Sp(1),
			Cell(running, 7), Sp(1),
			*ws.Name,
		)
	}
}
