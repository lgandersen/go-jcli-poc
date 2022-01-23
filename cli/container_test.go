package cli

import (
	"bytes"
	"fmt"
	"io"
	Openapi "jcli/client"
	"os"
	"os/exec"
	"testing"

	"gotest.tools/assert"
)

func TestContainerCreateListRemove(t *testing.T) {
	var container_id string
	var all = true
	t.Run("list container expecting empty list", ContainerListExpectEmptyListing(all))
	t.Run("create a new container", SuccesfullyCreateContainer(t, &container_id, "testerer", "base", []string{"/bin/ls"}))
	t.Run("verify that the container is NOT being listed, with all set to false", ContainerListExpectEmptyListing(!all))
	t.Run("verify that the container is being listed, with all set to true", ContainerListExpectContainer(all, "testerer"))
	t.Run("remove container", SuccesfullyRemoveContainer(t, container_id))
	t.Run("list container expecting empty list again", ContainerListExpectEmptyListing(all))
}

func TestContainerStartingAndStopping(t *testing.T) {
	var container_id string
	var attach = false
	var list_all = false
	t.Run("create  container that sleeps when started", SuccesfullyCreateContainer(t, &container_id, "testerer", "base", []string{"/bin/sleep", "10"}))
	t.Run("start container", StartContainer("testerer", attach))
	t.Run("verify that the container is running", VerifyRunningContainer(container_id))
	t.Run("check that running container is listed when 'all' is false", ContainerListExpectContainer(list_all, "testerer"))
	t.Run("stop container", StopContainer("testerer"))
	t.Run("verify container is stopped", VerifyStoppedContainer(container_id))
	t.Run("remove container", SuccesfullyRemoveContainer(t, container_id))
}

func TestContainerAttachAndStart(t *testing.T) {
	var container_id string
	var attach = true
	t.Run("create a new container", SuccesfullyCreateContainer(t, &container_id, "testerer", "base", []string{"/bin/ls"}))
	stdout := RunCommandCollectStdOut(func() { StartContainer("testerer", attach)(t) })
	expected_stdout := fmt.Sprintf(".cshrc\n.profile\nCOPYRIGHT\nbin\nboot\ndev\netc\nlib\nlibexec\nmedia\nmnt\nnet\nproc\nrescue\nroot\nsbin\nsys\ntmp\nusr\nvar\ncontainer %s stopped\n", container_id)
	assert.Equal(t, stdout, expected_stdout)
	t.Run("remove container", SuccesfullyRemoveContainer(t, container_id))
}

func TestContainerReferencingFunctionality(t *testing.T) {
	// FIXME: "try stopping a container that is already stopped"
	var container_id string
	var attach = false
	t.Run("create a container that sleeps when started", SuccesfullyCreateContainer(t, &container_id, "testerer", "base", []string{"/bin/sleep", "10"}))
	t.Run("start container using container id", StartContainer(container_id, attach))
	t.Run("stop container using container id", StopContainer(container_id))
	t.Run("start container using first part of the container id", StartContainer(container_id[:8], attach))
	t.Run("stop container using first part of the container id", StopContainer(container_id[:8]))
	t.Run("removed container", SuccesfullyRemoveContainer(t, container_id))
}

func StopContainer(container_id string) func(t *testing.T) {
	return func(t *testing.T) {
		response, err := ContainerStop([]string{container_id})
		assert.NilError(t, err)
		assert.Equal(t, len(response.JSON200.Id), 12)
	}
}

func StartContainer(container_id string, attach bool) func(*testing.T) {
	return func(t *testing.T) {
		if attach {
			StartAndAttachToContainer([]string{container_id})
		} else {
			container_ids := StartSeveralContainers([]string{container_id})
			container_id_returned := container_ids[0]
			assert.Equal(t, len(container_id_returned), 12)
		}
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

func VerifyRunningContainer(container_id string) func(t *testing.T) {
	return func(t *testing.T) {
		cmd := exec.Command("/bin/sh", "-c", "jls | grep "+container_id)
		err := cmd.Run()
		// If the container exists, grepping after the id will result in non-empty output from grep
		// which in turn results in exitcode 0 (non-zero if empty result from grep)
		assert.NilError(t, err)
	}
}

func ContainerListExpectEmptyListing(all bool) func(*testing.T) {
	return func(t *testing.T) {
		response, err := GetContainerList(all)
		assert.NilError(t, err)
		assert.Assert(t, response.JSON200 != nil)
		assert.Equal(t, len(*response.JSON200), 0)
	}
}

func ContainerListExpectContainer(all bool, container_name string) func(*testing.T) {
	return func(t *testing.T) {
		response, err := GetContainerList(all)
		assert.NilError(t, err)
		assert.Assert(t, response.JSON200 != nil)
		assert.Equal(t, len(*response.JSON200), 1)
		container_list := *response.JSON200
		assert.Equal(t, *container_list[0].Name, container_name)
	}
}

func SuccesfullyRemoveContainer(t *testing.T, container_id string) func(*testing.T) {
	return func(t *testing.T) {
		args := []string{container_id}
		response, err := PostContainerRemove(args)
		assert.NilError(t, err)
		var empty_id_response *Openapi.IdResponse
		assert.Assert(t, empty_id_response != response.JSON200)
		id := (*response.JSON200).Id
		assert.Equal(t, id, container_id)
	}
}

func SuccesfullyCreateContainer(t *testing.T, container_id *string, name, image string, cmd []string) func(t *testing.T) {
	config := Openapi.ContainerCreateJSONRequestBody{
		Networks:  &([]string{}),
		Volumes:   &([]string{}),
		Env:       &([]string{}),
		JailParam: &([]string{}),
	}
	return SuccesfullyCreateCustomContainer(t, container_id, name, image, cmd, config)
}

func SuccesfullyCreateCustomContainer(t *testing.T, container_id *string, name, image string, cmd []string, config Openapi.ContainerCreateJSONRequestBody) func(t *testing.T) {
	return func(t *testing.T) {
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

func StartContainerCollectOutput(t *testing.T, container_id string) string {
	attach := true
	stdout := RunCommandCollectStdOut(func() { StartContainer(container_id, attach)(t) })
	return stdout
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
