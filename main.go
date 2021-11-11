package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

var mdtext []string

func pprint(md string, lightMode *bool) {
	var mode = "dark"
	if *lightMode {
		mode = "light"
	}
	out, err := glamour.Render(md, mode)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Print(out)
}

func main() {

	var modeFlag = flag.Bool("l", false, "light mode (dark mode is default)")
	var fromFile = flag.String("f", "", "from file filename")
	flag.Parse()

	// https://stackoverflow.com/a/43947435/14494128
	fd1, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal("Something went wrong!\n%s\n", err)
	}
	if fd1.Mode()&os.ModeCharDevice == 0 {
		// from pipe
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			mdtext = append(mdtext, scanner.Text()+"\n")
		}
		if scanner.Err() != nil {
			log.Fatal("Something went wrong!\n%s\n", scanner.Err())
		}

	} else {
		// from file
		fmt.Println(*fromFile)
		data, err := os.Open(*fromFile)
		if err != nil {
			log.Fatal("Something went wrong!\n%s\n", err)
		}
		scanner := bufio.NewScanner(data)
		for scanner.Scan() {
			mdtext = append(mdtext, scanner.Text()+"\n")
		}
		if scanner.Err() != nil {
			log.Fatal("Something went wrong!\n%s\n", scanner.Err())
		}

	}

	pprint(strings.Join(mdtext, " "), modeFlag)
}