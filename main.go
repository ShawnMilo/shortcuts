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
		if strings.HasPrefix(line, "nnf") {
			printNotNilFatal()
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
		} else if strings.HasPrefix(line, "clog(") {
			consoleLog(t)
		} else if strings.HasPrefix(line, "clogVar") {
			consoleLogVar(line)
		} else if line == "gomain" {
			goMain()
		} else if line == "serve" {
			goWebserver()
		} else if line == "pymain" {
			pyMain()
		} else if line == "ul" {
			unorderedList(len(t) - len(strings.TrimLeft(line, " \t")))
		} else if line == "html5" {
			html5()
		} else if line == "ow:" {
			pyOpenWrite(t)
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

func printNotNilFatal() {
	// print error block
	fmt.Println("if err != nil{")
	fmt.Println(`log.Fatalf("Failed to do something: %s\n", err)`)
	fmt.Println("}")
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

func consoleLog(line string) {
	fmt.Println(strings.Replace(line, "clog(", "console.log(", 1))
}

func consoleLogVar(line string) {
	name := strings.Split(line, " ")[1]
	fmt.Printf("console.Log(\"%s: \" + %s);\n", name, name)
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

func goWebserver() {
	fmt.Println(`package main

import (
    "fmt"
	"net/http"
)

func main() {
    http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "hello")
}
`)
}

func pyMain() {
	fmt.Println(`#!/usr/bin/env python
"""
You should probably write something here.
"""

from __future__ import unicode_literals

def main():
    """
    Do the thing.
    """
    print "python"

if __name__ == '__main__':
    main()
`)
}

func html5() {
	fmt.Println(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>title</title>
		<link rel="stylesheet" href="./css/style.css" type="text/css">
		<meta name="viewport" content="width-device-width, initial-scale=1">
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js"></script>
	</head>
	<body>
		<div>
			<p>content</p>
		</div>
	</body>
</html>
`)
}

func httpHandlerFunc(line string) {
	name := strings.Split(line, " ")[1]
	fmt.Printf("func %s(w http.ResponseWriter, r *http.Request){\n", name)
}

func pyOpenWrite(line string) {
	l := len(line) - len(strings.TrimLeft(line, " "))
	pad := strings.Repeat(" ", l)
	lines := []string{
		`with open("out.txt", "wb") as raw:`,
		`    raw.write("{0}\n".format(msg))`,
	}

	for _, line = range lines {
		fmt.Printf("%s%s\n", pad, line)
	}

}

func unorderedList(margin int) {
	padding := strings.Repeat(" ", margin)
	fmt.Printf(padding)
	fmt.Println("<ul>")
	for i := 0; i < 3; i++ {
		fmt.Printf(padding)
		fmt.Println("\t<li>")
		fmt.Printf(padding)
		fmt.Println("\t\tthing")
		fmt.Printf(padding)
		fmt.Println("\t</li>")
	}
	fmt.Printf(padding)
	fmt.Println("</ul>")

}
