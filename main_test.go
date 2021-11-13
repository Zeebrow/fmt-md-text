package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestInputFromPipe(t *testing.T) {
	// note: need to test error codes using c.Wait()
	// os independent way to grab exit code?
	// https://stackoverflow.com/a/10385867/14494128

	c := exec.Command("./fmt-md-text")
	var stdout, stderr bytes.Buffer
	// var testText string = "asdfasdfsdf"
	const testText string = "# mdheading"

	const testOutputColorDarkMode string = "\033[38;5;228;48;5;63;1m"
	// const testOutputColorDarkMode string = "asdfasdfasdf"
	const testOutput string = "mdheading"
	// const testOutput string = "asdfasdfsdf"
	const testOutputReset string = "[0m\033[38;5;228;48;5;63;1m\033[0m"
	// const testOutputReset string = "adsfkhaksdhfklashdfhsdf"

	c.Stdin = strings.NewReader(testText)
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()

	var fd1 string = stdout.String()
	var fd2 string = stderr.String()

	t.Logf("stdout: %v", fd1)
	t.Logf("stderr: %v", fd2)
	if err != nil {
		t.Errorf("Error should be nil for file piped to fmt-md-text (how in the world..?). \nstdout: %s\nsterr: %v", fd1, fd2)
	}

	if !(strings.Contains(fd1, testOutputColorDarkMode)) {
		t.Error("'dark mode' output expected (-l flag is absent)")
	}
	if !(strings.Contains(fd1, testOutput)) {
		t.Error("Output should be pretty much the same as input")
		t.Error(fd1)
	}
	if !(strings.Contains(fd1, testOutputReset)) {
		t.Error("Colored output is never reset")
	}

}

func TestFileAsArg(t *testing.T) {
	c := exec.Command("./fmt-md-text", "-f", "test-fixtures/example.md")
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		t.Error("Errot should be nil for file that exists")
	}
}
func TestNoSuchFile(t *testing.T) {
	c := exec.Command("./fmt-md-text", "-f", "asdf")
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if err == nil {
		t.Error("Errot should not be nil for file that does not exist")
	}
}
