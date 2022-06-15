package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestNoSuchFileLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.SkipNow()
		return
	}
	var expectedErrorMsg string
	switch runtime.GOOS {
	case "windows":
		expectedErrorMsg = "open asdf: The system cannot find the file specified."
	case "linux":
		expectedErrorMsg = "open asdf: no such file or directory"
	default:
		fmt.Errorf("Unsupported platform '%s'\n", runtime.GOOS)
	}
	var b Executable
	b.SetTestBinaryName()
	c := exec.Command(b.fullname, "-f", "asdf")

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	fd1 := fmt.Sprint(c.Stdout)
	// if fd1 != expectedErrorMsg {
	if !strings.Contains(fd1, expectedErrorMsg) {
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
