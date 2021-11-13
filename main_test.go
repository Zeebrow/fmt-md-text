package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestFileFromPipe(t *testing.T) {
	// c := exec.Command("cat", "test-fixtures/example.md", "|", "./fmt-md-text")
	c := exec.Command("./fmt-md-text")
	var stdin, stdout, stderr bytes.Buffer

	c.Stdin = &stdin
	c.Stdin = strings.NewReader("# mdheading\n\nmdtexts")
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		t.Error("Errot should be nil for file piped to fmt-md-text")
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
