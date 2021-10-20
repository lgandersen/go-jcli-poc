package cli

import (
	"bytes"
	"io"
	Openapi "jcli/client"
	"os"
	"os/exec"
	"testing"

	"gotest.tools/assert"
)

var testerer_container_id *string

func TestContainerCreateListRemove(t *testing.T) {
	var container_id string
	t.Run("list container expecting empty list", ContainerListExpectEmptyListing)
	t.Run("create a new container", CreateContainer(t, &container_id, "testerer", "base", []string{"/bin/ls"}))
	t.Run("verify that the container is being listed", ContainerListExpectContainer(t, "testerer"))
	t.Run("removed container", ContainerRemoveTesterer(t, container_id))
	t.Run("list container expecting empty list again", ContainerListExpectEmptyListing)
}

func TestContainerCreateStartRemove(t *testing.T) {
	var container_id string
	t.Run("create container that sleeps when started", CreateContainer(t, &container_id, "testerer", "base", []string{"/bin/sleep", "10"}))
	t.Run("start container", StartContainer("testerer"))
	VerifyRunningContainer(t, container_id)
	t.Run("stop container", StopContainer("testerer"))
	t.Run("verify container is stopped", VerifyStoppedContainer(container_id))
	t.Run("removed container", ContainerRemoveTesterer(t, container_id))
}

func StopContainer(container_id string) func(t *testing.T) {
	return func(t *testing.T) {
		response, err := ContainerStop([]string{container_id})
		assert.NilError(t, err)
		assert.Equal(t, len(response.JSON200.Id), 12)
	}
}

func StartContainer(container_id string) func(*testing.T) {
	return func(t *testing.T) {
		container_ids := StartSeveralContainers([]string{container_id})
		container_id_returned := container_ids[0]
		assert.Equal(t, len(container_id_returned), 12)
	}
}

func VerifyStoppedContainer(container_id string) func(*testing.T) {
	return func(t *testing.T) {
		cmd := exec.Command("/bin/sh", "-c", "jls | grep "+container_id)
		cmd.Run()
		output, _ := cmd.Output()
		assert.Equal(t, len(output), 0)
	}
}

func VerifyRunningContainer(t *testing.T, container_id string) {
	cmd := exec.Command("/bin/sh", "-c", "jls | grep "+container_id)
	err := cmd.Run()
	// If the container exists, grepping after the id will result in non-empty output from grep
	// which in turn results in exitcode 0 (non-zero if empty result from grep)
	assert.NilError(t, err)
}

func ContainerListExpectEmptyListing(t *testing.T) {
	all := true
	response, err := GetContainerList(all)
	assert.NilError(t, err)
	assert.Assert(t, response.JSON200 != nil)
	assert.Equal(t, len(*response.JSON200), 0)
}

func ContainerListExpectContainer(t *testing.T, container_name string) func(*testing.T) {
	return func(t *testing.T) {
		all := true
		response, err := GetContainerList(all)
		assert.NilError(t, err)
		assert.Assert(t, response.JSON200 != nil)
		assert.Equal(t, len(*response.JSON200), 1)
		container_list := *response.JSON200
		assert.Equal(t, *container_list[0].Name, container_name)
	}
}

func ContainerRemoveTesterer(t *testing.T, container_id string) func(*testing.T) {
	return func(t *testing.T) {
		args := []string{container_id}
		response, err := PostContainerRemove(args)
		assert.NilError(t, err)
		var empty_id_response *Openapi.IdResponse
		assert.Assert(t, empty_id_response != response.JSON200)
		var id string
		id = (*response.JSON200).Id
		assert.Equal(t, id, container_id)
	}
}

func CreateContainer(t *testing.T, container_id *string, name, image string, cmd []string) func(t *testing.T) {
	return func(t *testing.T) {
		config := Openapi.ContainerCreateJSONRequestBody{
			Networks:  &([]string{}),
			Volumes:   &([]string{}),
			Env:       &([]string{}),
			JailParam: &([]string{}),
		}
		var args = make([]string, len(cmd)+1)
		args[0] = image
		copy(args[1:], cmd[:])
		response, err := PostContainerCreate(&name, config, args)
		assert.NilError(t, err)
		var empty_json201 *Openapi.IdResponse
		assert.Assert(t, response.JSON201 != empty_json201)
		assert.Equal(t, len(response.JSON201.Id), 12)
		*container_id = response.JSON201.Id
	}
}

func RunCommandCollectStdOut(f func()) string {
	old_stdout := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	outC := make(chan string)

	f()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old_stdout
	output := <-outC
	return string(output)
}
