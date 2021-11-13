package cli

import (
	"context"
	"errors"
	"fmt"
	Openapi "jcli/client"

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

func NetworkListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "list",
		Short:                 "List networks",
		Long:                  `List networks`,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := NetworkList()
			if err == nil {
				PrintNetworkList(response.JSON200)
			}
		},
	}

	return cmd
}
func NetworkList() (*Openapi.NetworkListResponse, error) {
	client := NewHTTPClient()
	response, err := client.NetworkListWithResponse(context.TODO())
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

func PrintNetworkList(networks *[]Openapi.NetworkSummary) {
	fmt.Println(Cell("NETWORK ID", 12), Cell("NAME", 25), "DRIVER")
	for _, network := range *networks {
		fmt.Println(Cell(*network.Id, 12), Cell(*network.Name, 25), *network.Driver)

	}
}

func NetworkCreateCommand() *cobra.Command {
	var driver, ifname, name, subnet string
	config := Openapi.NetworkCreateJSONRequestBody{
		Driver: &driver,
		Ifname: &ifname,
		Name:   name,
		Subnet: &subnet,
	}

	cmd := &cobra.Command{
		Use:                   "create [OPTIONS] NETWORK_NAME",
		Short:                 "Create a new network",
		Long:                  `Create a new network`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run:                   func(cmd *cobra.Command, args []string) { NetworkCreate(args, config) },
	}

	flags := cmd.Flags()
	flags.StringVarP(config.Driver, "driver", "d", "", "Which driver to use for the network. Only 'loopback' is possible atm.")
	flags.StringVar(config.Ifname, "ifname", "", "Name of the loopback interface used for the network")
	flags.StringVar(config.Subnet, "subnet", "", "Subnet in CIDR format that represents the network segment")
	return cmd
}
func NetworkCreate(args []string, config Openapi.NetworkCreateJSONRequestBody) {
	config.Name = args[0]
	client := NewHTTPClient()
	response, err := client.NetworkCreateWithResponse(context.TODO(), config)
	err = verify_response(response.StatusCode(), 201, err)
	if err == nil {
		fmt.Println(response.JSON201.Id)
	}
}

func NetworkRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "remove NETWORK",
		Short:                 "Remove a network",
		Long:                  `Remove a network`,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run:                   func(cmd *cobra.Command, args []string) { RemoveNetworks(args) },
	}
	return cmd
}

func RemoveNetworks(name_or_ids []string) {
	for _, name_or_id := range name_or_ids {
		client := NewHTTPClient()
		response, err := client.NetworkRemoveWithResponse(context.TODO(), name_or_id)
		err = verify_response(response.StatusCode(), 200, err)
		if err == nil {
			fmt.Println(response.JSON200.Id)
		}
	}
}

func NetworkConnectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "connect NETWORK CONTAINER",
		Short:                 "Connect a container to a network",
		Long:                  `Connect a container to a network`,
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Run:                   func(cmd *cobra.Command, args []string) { NetworkConnect(args) },
	}
	return cmd
}

func NetworkConnect(args []string) {
	network_name := args[0]
	container_name := args[1]
	client := NewHTTPClient()
	response, err := client.NetworkConnectWithResponse(context.TODO(), network_name, container_name)
	err = verify_response(response.StatusCode(), 204, err)
	if err == nil {
		fmt.Println(response.JSON204.Id)
	}
}

func NetworkDisconnectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "disconnect NETWORK CONTAINER",
		Short:                 "Disconnect a container from a network",
		Long:                  `Disconnect a container from a network`,
		Args:                  cobra.ExactArgs(2),
		DisableFlagsInUseLine: true,
		Run:                   func(cmd *cobra.Command, args []string) { NetworkConnect(args) },
	}
	return cmd
}

func NetworkDisconnect(args []string) {
	network_name := args[0]
	container_name := args[1]
	client := NewHTTPClient()
	response, err := client.NetworkDisconnectWithResponse(context.TODO(), network_name, container_name)
	err = verify_response(response.StatusCode(), 204, err)
	if err == nil {
		fmt.Println(response.JSON204.Id)
	}
}
