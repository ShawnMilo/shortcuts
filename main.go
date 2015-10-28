package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if stat.Mode()&os.ModeCharDevice != 0 {
		log.Fatal("please pipe in some data")
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := s.Text()
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "nnf:") {
			printNotNilFatal(t)
		} else if strings.HasPrefix(line, "lpf(") {
			printLogPrintf(t)
		} else if strings.HasPrefix(line, "fpf(") {
			fPrintf(t)
		} else if strings.HasPrefix(line, "fpl(") {
			fPrintln(t)
		} else if strings.HasPrefix(line, "hfunc") {
			httpHandlerFunc(t)
		} else if strings.HasPrefix(line, "lpl(") {
			lPrintln(t)
		} else if line == "gomain" {
			goMain()
		} else if line == "pymain" {
			pyMain()
		} else if line == "ubb" {
			fmt.Println("#!/usr/bin/env bash")
		} else if line == "ubp" {
			fmt.Println("#!/usr/bin/env python")
		} else if strings.Contains(t, "(t.t)") {
			fmt.Println(strings.Replace(t, "(t.t)", "(t *testing.T)", 1))
		} else {
			fmt.Println(t)
		}
	}
}

func printNotNilFatal(line string) {
	l := len(line) - len(strings.TrimLeft(line, " "))
	pad := strings.Repeat(" ", l)
	msg := strings.SplitN(line, ":", 2)[1]

	// print error block
	fmt.Printf("%sif err != nil{\n", pad)
	fmt.Printf("%s%slog.Fatal(%q, err)\n", pad, pad, msg+": ")
	fmt.Printf("%s}\n", pad)
}

func printLogPrintf(line string) {
	fmt.Println(strings.Replace(line, "lpf(", "log.Printf(", 1))
}

func fPrintf(line string) {
	fmt.Println(strings.Replace(line, "fpf(", "fmt.Printf(", 1))
}

func lPrintln(line string) {
	fmt.Println(strings.Replace(line, "lpl(", "log.Println(", 1))
}

func fPrintln(line string) {
	fmt.Println(strings.Replace(line, "fpl(", "fmt.Println(", 1))
}

func goMain() {
	fmt.Println(`package main

import (
    "fmt"
)

func main() {
    fmt.Println("gopher")
}
`)
}

func pyMain() {
	fmt.Println(`#!/usr/bin/env python
"""
You should probably write something here.
"""

from __future__ import unicode_literals

func main():
    """
    Do the thing.
    """
    print "python"

if __name__ == '__main__':
    main()
`)
}

func httpHandlerFunc(line string) {
	name := strings.Split(line, " ")[1]
	fmt.Printf("func %s(w http.ResponseWriter, r *http.Request){\n", name)
}
