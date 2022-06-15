package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

type Executable struct {
	prefix   string
	name     string
	suffix   string
	fullname string
}

func (e *Executable) SetTestBinaryName() {
	switch runtime.GOOS {
	case "windows":
		e.prefix = ".\\"
		e.suffix = ".exe"
	case "linux":
		e.prefix = "./"
		e.suffix = ""
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		os.Exit(1)
	}
	binName := os.Getenv("FMT_MD_TEXT_BINARY")
	if binName == "" {
		e.name = "fmt-md-text"
	} else {
		e.name = binName
	}
	e.fullname = e.prefix + e.name + e.suffix
	return
}

func TestVersionString(t *testing.T) {
	var b Executable
	b.SetTestBinaryName()
	var stdout, stderr bytes.Buffer
	c := exec.Command(b.fullname, "-version")
	c.Stdout = &stdout
	c.Stderr = &stderr
	err := c.Run()
	fd1 := fmt.Sprint(c.Stdout)
	if err != nil {
		t.Errorf("Unknown error\n")
	}
	if b.name != "fmt-md-text" {
		if strings.HasPrefix(fd1, "dev") {
			t.Errorf("builds compiled with `make` should not have a 'dev' version string (got %s)\n", fd1)
		}
	} else {
		fmt.Println("got a quickone")
		if !strings.HasPrefix(fd1, "dev ") {
			t.Errorf("Expected 'dev' prefix to version string (string was '%s')\n", fd1)
		}
	}
}

// const testOutputColorLightMode = "\033[38;5;203;48;5;254m"
// const testBacktickOutputDarkMode = "\033[38;5;203;48;5;236m"

func TestInputFromPipeLight(t *testing.T) {
	// note: need to test error codes using c.Wait()
	// os independent way to grab exit code?
	// https://stackoverflow.com/a/10385867/14494128
	var b Executable
	b.SetTestBinaryName()
	c := exec.Command(b.fullname, "-l")
	var stdout, stderr bytes.Buffer
	// var testText string = "asdfasdfsdf"
	const testText string = "`mdcodelight`"

	// const testOutputColorDarkMode string = "\033[38;5;228;48;5;63;1m"
	// const testOutputColorDarkMode string = "asdfasdfasdf"
	// const testOutput string = "mdheading"
	// const testOutput string = "asdfasdfsdf"
	// const testOutputReset string = "[0m\033[38;5;228;48;5;63;1m\033[0m"
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

	if !(strings.Contains(fd1, "\033[38;5;203;48;5;254m")) {
		t.Error("'dark mode' output expected (-l flag is absent)")
	}
	if !(strings.Contains(fd1, "mdcodelight")) {
		t.Error("Output should be pretty much the same as input")
		t.Error(fd1)
	}
	if !(strings.Contains(fd1, "\033[0m")) {
		t.Error("Colored output is never reset")
	}

}
func TestInputFromPipeDark(t *testing.T) {
	var b Executable
	b.SetTestBinaryName()
	c := exec.Command(b.fullname)
	var stdout, stderr bytes.Buffer
	const testText string = "`mdcode`\n\n"

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

	if !(strings.Contains(fd1, "\033[38;5;203;48;5;236m")) {
		t.Error("'dark mode' output expected (-l flag is absent)")
	}
	if !(strings.Contains(fd1, "mdcode")) {
		t.Error("Output should be pretty much the same as input")
		t.Error(fd1)
	}
	if !(strings.Contains(fd1, "\033[0m")) {
		t.Error("Colored output is never reset")
	}
}

func TestNoSuchFile(t *testing.T) {
	//@@@
	t.SkipNow()
	var b Executable
	b.SetTestBinaryName()
	expectedErrorMsg := "open asdf: no such file or directory\n"
	c := exec.Command(b.fullname, "-f", "asdf")

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
