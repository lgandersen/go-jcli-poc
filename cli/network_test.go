package cli

import (
	"encoding/json"
	"fmt"
	Openapi "jcli/client"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNetworkListing(t *testing.T) {
	var create_id, remove_id string
	verify_ids := func(t *testing.T) { assert.Equal(t, create_id, remove_id) }
	t.Run("verify that the list is empty in the beginning", ListNetworksExpectIds(t, []string{"default", "host"}))
	t.Run("create a new network", SuccesfullyCreateNetwork(t, &create_id, "testnet"))
	t.Run("verify that the newly created network shows up in the list", ListNetworksExpectIds(t, []string{"default", "host", "testnet"}))
	t.Run("remove the newly created network", SuccesfullyRemoveANetwork(t, &remove_id, "testnet"))
	t.Run("verify that the network is not being listed anymore", ListNetworksExpectIds(t, []string{"default", "host"}))
	t.Run("verify that network-create and network-remove returns the same id", verify_ids)
}

func TestConnectContainerToCustomNetworkAtCreationTime(t *testing.T) {
	var container_id, network_id string
	expected_if_names := []string{"jclitest"}
	config := Openapi.ContainerCreateJSONRequestBody{
		Networks:  &([]string{"testnet"}),
		Volumes:   &([]string{}),
		Env:       &([]string{}),
		JailParam: &([]string{}),
	}
	cmd := []string{"netstat", "--libxo", "json", "-4", "-i"}

	t.Run("create network 'testnet' to use in testing", SuccesfullyCreateNetwork(t, &network_id, "testnet"))
	t.Run("create container that is attached to 'testnet'", SuccesfullyCreateCustomContainer(t, &container_id, "nettester", "base", cmd, config))
	t.Run("Start container and verify that it is only connected to 'testnet'", StartContainerAndVerifyInterfaces(t, &container_id, expected_if_names))
	t.Run("Disconnect the container from 'testnet'", SuccesfullyDisconnectContainerToNetwork(t, network_id, container_id))
	t.Run("Start container and verify that it is not connected to anything", StartContainerAndVerifyInterfaces(t, &container_id, []string{}))
	t.Run("remove container", SuccesfullyRemoveContainer(t, container_id))
	t.Run("destroy network after use", SuccesfullyRemoveANetwork(t, &network_id, "testnet"))
}

func TestConnectContainerToCustomNetworkAfterCreatedWithDefaultNetwork(t *testing.T) {
	var container_id, network_id string
	expected_if_names := []string{"jocker0", "jclitest"}
	config := Openapi.ContainerCreateJSONRequestBody{
		Networks:  &([]string{}),
		Volumes:   &([]string{}),
		Env:       &([]string{}),
		JailParam: &([]string{}),
	}
	cmd := []string{"netstat", "--libxo", "json", "-4", "-i"}

	t.Run("create network 'testnet' to use in testing", SuccesfullyCreateNetwork(t, &network_id, "testnet"))
	t.Run("create container that is attached to the default network", SuccesfullyCreateCustomContainer(t, &container_id, "nettester", "base", cmd, config))
	t.Run("connect container to 'testnet' network", SuccesfullyConnectContainerToNetwork(t, network_id, container_id))
	t.Run("Start container and verify that both networks are connected", StartContainerAndVerifyInterfaces(t, &container_id, expected_if_names))
	t.Run("Disconnect the container from 'testnet'", SuccesfullyDisconnectContainerToNetwork(t, network_id, container_id))
	t.Run("Start container and verify that the container is only connected to the default network", StartContainerAndVerifyInterfaces(t, &container_id, []string{"jocker0"}))
	t.Run("remove container", SuccesfullyRemoveContainer(t, container_id))
	t.Run("destroy network after use", SuccesfullyRemoveANetwork(t, &network_id, "testnet"))
}

func SuccesfullyConnectContainerToNetwork(t *testing.T, network_id, container_id string) func(*testing.T) {
	return func(t *testing.T) {
		_, err := NetworkConnect([]string{network_id, container_id})
		assert.NilError(t, err)
	}
}

func SuccesfullyDisconnectContainerToNetwork(t *testing.T, network_id, container_id string) func(*testing.T) {
	return func(t *testing.T) {
		_, err := NetworkDisconnect([]string{network_id, container_id})
		assert.NilError(t, err)
	}
}

func StartContainerAndVerifyInterfaces(t *testing.T, container_id *string, expected_if_names []string) func(*testing.T) {
	return func(t *testing.T) {
		stdout := StartContainerCollectOutput(t, *container_id)
		netstat_if_status := DecodeNetstatInterfaceStatus(*container_id, stdout)
		if_names := ExtractInterfacesFromIfStatus(netstat_if_status)
		assert.DeepEqual(t, if_names, expected_if_names)
	}
}

func ExtractInterfacesFromIfStatus(status NetstatInterfaceStatus) []string {
	interface_names := make([]string, len(status.Statistics.Interface))
	for idx, interface_info := range status.Statistics.Interface {
		interface_names[idx] = interface_info.Name
	}
	return interface_names
}

func DecodeNetstatInterfaceStatus(container_id, output string) NetstatInterfaceStatus {
	var netstat_status NetstatInterfaceStatus
	ending_msg := fmt.Sprintf("container %s stopped", container_id)
	output_json := strings.Replace(output, ending_msg, "", -1)
	err := json.Unmarshal([]byte(output_json), &netstat_status)
	if err != nil {
		fmt.Println("error decoding netstat json output:", err)
	}
	return netstat_status
}

type NetstatInterfaceStatus struct {
	Statistics struct {
		Interface []InterfaceInfo `json:"interface"`
	} `json:"statistics"`
}

type InterfaceInfo struct {
	Name            string `json:"name"`
	Flags           string `json:"flags"`
	Network         string `json:"network"`
	Address         string `json:"address"`
	ReceivedPackets int    `json:"received-packets"`
	SentPackets     int    `json:"sent-packets"`
}

func SuccesfullyRemoveANetwork(t *testing.T, network_id *string, name_or_id string) func(*testing.T) {
	return func(t *testing.T) {
		response, errs := RemoveNetworks([]string{name_or_id})
		assert.NilError(t, errs[0])
		*network_id = response[0].JSON200.Id
	}
}

func SuccesfullyCreateNetwork(t *testing.T, network_id *string, network_name string) func(*testing.T) {
	return func(t *testing.T) {
		driver := "loopback"
		ifname := "jclitest"
		subnet := "10.13.37.0/24"
		config := Openapi.NetworkCreateJSONRequestBody{
			Driver: &driver,
			Ifname: &ifname,
			Name:   network_name,
			Subnet: &subnet,
		}
		response, err := NetworkCreate([]string{network_name}, config)
		assert.NilError(t, err)

		fmt.Println("WHAAAAAT", response.JSON201.Id)
		*network_id = response.JSON201.Id
	}
}

func ListNetworksExpectIds(t *testing.T, expected_names []string) func(*testing.T) {
	return func(t *testing.T) {
		response, err := NetworkList()
		assert.NilError(t, err)
		assert.Equal(t, len(expected_names), len(*response.JSON200))
		for idx, expected_id := range expected_names {
			received_name := *(*response.JSON200)[idx].Name
			assert.Equal(t, received_name, expected_id)
		}
	}
}
