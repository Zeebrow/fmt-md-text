package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"testing"
)

func TestNoSuchFileLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}
	expectedErrorMsg := "open asdf: no such file or directory\n"
	c := exec.Command(GetTestBinaryName(), "-f", "asdf")

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	fd1 := fmt.Sprint(c.Stdout)
	if fd1 != expectedErrorMsg {
		t.Error("stdout message does not match expceted response.")
		t.Errorf("Got: %s", fd1)
		t.Errorf("Expected: %s", expectedErrorMsg)
	}
	t.Logf("stdout: %v", c.Stdout)
	t.Logf("stderr: %v", c.Stderr)
	if err == nil {
		t.Error(err)
		t.Error("Error should not be nil for file that does not exist")
		t.Errorf("stdout: %v\n", c.Stdout)
		t.Errorf("stderr: %v\n", c.Stderr)
	}
}
