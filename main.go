/*
Copyright (C) <year>  <name of author>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// usr/bin/env go run "$0" "$@";
// above is about equal to a bash shebang
// you can ./main.go -f example.md and its kinda close enough
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
)

type ProgramInfo struct {
	progName   string
	version    string
	buildDate  string
	commitHash string
}

func (pi *ProgramInfo) show() {
	fmt.Printf("%s version '%s' (%s) built on '%s'\n", pi.progName, pi.version, pi.commitHash, pi.buildDate)
}

var ProgramName = "drcat"
var Version = "dev"
var BuildDate = ""
var CommitHash = ""
var ProgInfo = ProgramInfo{
	progName:   ProgramName,
	version:    Version,
	buildDate:  BuildDate,
	commitHash: CommitHash,
}

// must be between 0-125 on uboontoo
const ERR_WTF = 66
const ERR_NOF = 69
const ERR_RTFM = 42

var mdtext []string

// func pprint(md string, lightMode *bool) {
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
		fmt.Println(err.Error())
		os.Exit(ERR_WTF)
	}
	fmt.Print(out)
}

func usage() {
	fmt.Printf("usage: %s [-l] [-f FILE.md]\n", ProgInfo.progName)
}

func main() {

	var helpFlag = flag.Bool("h", false, "print usage")
	var modeFlag = flag.Bool("l", false, "light mode (dark mode is default)")
	var fromFile = flag.String("f", "", "from file filename")
	var versionFlag = flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *versionFlag {
		ProgInfo.show()
		os.Exit(0)
	}

	if *helpFlag {
		usage()
		os.Exit(0)
	}
	// https://stackoverflow.com/a/43947435/14494128
	fd1, err := os.Stdin.Stat()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(ERR_WTF)
	}
	if fd1.Mode()&os.ModeCharDevice == 0 {
		// from pipe
		pprint(os.Stdin, modeFlag)

	} else {
		// from file path, passed with -f
		if *fromFile == "" {
			fmt.Printf("No file provided! %s\n", *fromFile)
			usage()
			os.Exit(ERR_RTFM)
		}
		fi, err := os.Open(*fromFile)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(ERR_NOF)
		}
		pprint(fi, modeFlag)
	}

}
