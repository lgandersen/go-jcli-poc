package cli

import (
	Openapi "jcli/client"
	"testing"

	"gotest.tools/assert"
)

func TestContainerListSubCommand(t *testing.T) {
	t.Run("list container expecting empty list", ContainerListExpectEmptyListing)
	t.Run("create container named testerer", ContainerAddTesterer)
	t.Run("list container expecting testerer container", ContainerListExpectTesterer)
	t.Run("remove container named testerer", ContainerRemoveTesterer)
	t.Run("list container expecting empty list", ContainerListExpectEmptyListing)
}

var testerer_container_id *string

func ContainerListExpectEmptyListing(t *testing.T) {
	all := true
	response, err := GetContainerList(all)
	assert.NilError(t, err)
	assert.Assert(t, response.JSON200 != nil)
	assert.Equal(t, len(*response.JSON200), 0)
}

func ContainerAddTesterer(t *testing.T) {
	config := Openapi.ContainerCreateJSONRequestBody{
		Networks:  &([]string{}),
		Volumes:   &([]string{}),
		Env:       &([]string{}),
		JailParam: &([]string{}),
	}
	name := "testerer"
	args := []string{"base", "/bin/ls"}
	response, err := PostContainerCreate(&name, config, args)
	assert.NilError(t, err)
	var empty_json201 *Openapi.IdResponse
	assert.Assert(t, response.JSON201 != empty_json201)
	assert.Equal(t, len(response.JSON201.Id), 12)
	testerer_container_id = &response.JSON201.Id
}

func ContainerListExpectTesterer(t *testing.T) {
	all := true
	response, err := GetContainerList(all)
	assert.NilError(t, err)
	assert.Assert(t, response.JSON200 != nil)
	assert.Equal(t, len(*response.JSON200), 1)
	container_list := *response.JSON200
	assert.Equal(t, *container_list[0].Name, "testerer")
}

func ContainerRemoveTesterer(t *testing.T) {
	args := []string{"testerer"}
	response, err := PostContainerRemove(args)
	assert.NilError(t, err)
	var empty_id_response *Openapi.IdResponse
	assert.Assert(t, empty_id_response != response.JSON200)
	var id string
	id = (*response.JSON200).Id
	assert.Equal(t, id, *testerer_container_id)
}
