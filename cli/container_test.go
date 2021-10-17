package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestInitializationAndJockerEngineStartup(t *testing.T) {
	RecreateJockerZroot(t)
	stdout, cmd := InitJockerEngine(t)
	ShutdownJockerEngine(cmd, stdout)
}

func TestContainerSubCommand(t *testing.T) {
	RecreateJockerZroot(t)
	stdout, cmd := InitJockerEngine(t)
	//ContainerListTest(t)
	t.Run("Testing list container", ContainerListTest)
	ShutdownJockerEngine(cmd, stdout)
}

func ContainerListTest(t *testing.T) {
	all := true
	args := []string{}
	cli_out := RunCommandCollectStdOut(func() { RunContainerList(all, args) })
	expected_out := "CONTAINER ID       IMAGE                 COMMAND                       CREATED                  STATUS        NAME\n"
	if cli_out != expected_out {
		t.Errorf("client output != expected output.\nCilent: %s\nExpected: %s", cli_out, expected_out)
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

	fmt.Println(" === BEGIN JCLI OUTPUT ===")
	fmt.Print(string(output))
	fmt.Println("\n === END JCLI OUTPUT ===")
	return string(output)
}

const jocker_zroot = "zroot/jocker"
const jocker_engine_path = "/home/vagrant/jocker/"

func RecreateJockerZroot(t *testing.T) {
	exec.Command("/bin/sh", "-c", "sudo /sbin/zfs destroy -rf "+jocker_zroot).Run()
	exec.Command("/bin/sh", "-c", "sudo /sbin/zfs create "+jocker_zroot).Run()
}

func InitJockerEngine(t *testing.T) (io.ReadCloser, *exec.Cmd) {
	cmd := exec.Command("/usr/local/bin/sudo", jocker_engine_path+"_build/dev/rel/jockerd/bin/jockerd", "start")
	stdout, err := cmd.StdoutPipe()
	assert.Assert(t, err == nil)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
		assert.Assert(t, false)
	}
	time.Sleep(3 * time.Second)
	return stdout, cmd
}

func ShutdownJockerEngine(cmd *exec.Cmd, jockerd_stdout io.ReadCloser) {
	pid := cmd.Process.Pid
	err := exec.Command("/bin/kill", fmt.Sprint(pid)).Run()
	if err != nil {
		log.Fatal(err)
	}
	output, err := io.ReadAll(jockerd_stdout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" === BEGIN JOCKER ENGINE OUTPUT ===")
	fmt.Print(string(output))
	fmt.Println("\n === END JOCKER ENGINE OUTPUT ===")
}
