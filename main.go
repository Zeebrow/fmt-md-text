package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

var ERR_WTF = 1
var ERR_NOF = 2
var ERR_RTFM = 3

var mdtext []string
var progname = "fmt-md-text"

func pprint(md string, lightMode *bool) {
	var mode = "dark"
	if *lightMode {
		mode = "light"
	}
	out, err := glamour.Render(md, mode)
	if err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Print(out)
}

func usage() {
	fmt.Printf("usage: %s [-l] [-f FILE.md]\n", progname)
}

func main() {

	var modeFlag = flag.Bool("l", false, "light mode (dark mode is default)")
	var fromFile = flag.String("f", "", "from file filename")
	flag.Parse()

	// https://stackoverflow.com/a/43947435/14494128
	fd1, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Something went wrong!\n%s\n", err)
		os.Exit(ERR_WTF)
	}
	if fd1.Mode()&os.ModeCharDevice == 0 {
		// from pipe
		fmt.Printf("pipe\n")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			mdtext = append(mdtext, scanner.Text()+"\n")
		}
		if scanner.Err() != nil {
			fmt.Printf("Something went wrong!\n%s\n", scanner.Err())
			os.Exit(ERR_WTF)
		}

	} else {
		// from file
		if *fromFile == "" {
			fmt.Printf("No file provided! %s\n", *fromFile)
			usage()
			os.Exit(ERR_RTFM)
		}
		if _, err := os.Stat(*fromFile); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("No such file: %s\n", *fromFile)
			os.Exit(ERR_NOF)
		}

		data, err := os.Open(*fromFile)
		if err != nil {
			fmt.Printf("Something went wrong!\n%s\n", err)
			os.Exit(3)
		}
		scanner := bufio.NewScanner(data)
		for scanner.Scan() {
			mdtext = append(mdtext, scanner.Text()+"\n")
		}
		if scanner.Err() != nil {
			fmt.Printf("Something went wrong!\n%s\n", scanner.Err())
			os.Exit(ERR_WTF)
		}

	}

	pprint(strings.Join(mdtext, " "), modeFlag)
}
