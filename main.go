package main

import (
	"bufio"
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

//func pprint(md string, lightMode *bool) {
func pprint(f *os.File, lightMode *bool) {
	var mode = "dark"
	if *lightMode {
		mode = "light"
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		mdtext = append(mdtext, scanner.Text()+"\n")
	}
	if scanner.Err() != nil {
		fmt.Printf("Something went wrong!\n%s\n", scanner.Err())
		os.Exit(ERR_WTF)
	}
	out, err := glamour.Render(strings.Join(mdtext, ""), mode)
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(ERR_WTF)
	}
	fmt.Print(out)
}

func usage() {
	fmt.Printf("usage: %s [-l] [-f FILE.md]\n", progname)
}

func main() {

	var helpFlag = flag.Bool("h", false, "print usage")
	var modeFlag = flag.Bool("l", false, "light mode (dark mode is default)")
	var fromFile = flag.String("f", "", "from file filename")
	flag.Parse()

	if *helpFlag {
		usage()
		os.Exit(0)
	}
	// https://stackoverflow.com/a/43947435/14494128
	fd1, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("Something went wrong!\n%s\n", err)
		os.Exit(ERR_WTF)
	}
	if fd1.Mode()&os.ModeCharDevice == 0 {
		// from pipe
		fmt.Printf("pipe\n")
		pprint(os.Stdin, modeFlag)

	} else {
		// from file
		if *fromFile == "" {
			fmt.Printf("No file provided! %s\n", *fromFile)
			usage()
			os.Exit(ERR_RTFM)
		}
		fi, err := os.Open(*fromFile)
		if err != nil {
			fmt.Printf("Something went wrong!\n%s\n", err)
			os.Exit(3)
		}
		pprint(fi, modeFlag)
	}

}
