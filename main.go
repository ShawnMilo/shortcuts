package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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
		line := s.Text()
		trim := strings.TrimSpace(s.Text())
		if strings.HasPrefix(trim, "nnf") {
			printNotNilFatal()
		} else if strings.HasPrefix(trim, "nnl") {
			printNotNilLog()
		} else if trim == "openFile" {
			openFile()
		} else if trim == "reqStdin" {
			reqStdin()
		} else if trim == "getURL" {
			getURL()
		} else if strings.HasPrefix(trim, "lpf(") {
			printLogPrintf(line)
		} else if strings.HasPrefix(trim, "fpf(") {
			fPrintf(line)
		} else if strings.HasPrefix(trim, "fpl(") {
			fPrintln(line)
		} else if strings.HasPrefix(trim, "hfunc") {
			httpHandlerFunc(line)
		} else if strings.HasPrefix(trim, "lpl(") {
			lPrintln(line)
		} else if trim == "gomain" {
			goMain()
		} else if trim == "tempFile" {
			tempFile()
		} else if trim == "now" {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		} else if trim == "serve" {
			goWebserver()
		} else if trim == "pymain" {
			pyMain()
		} else if trim == "ul" {
			unorderedList(len(line) - len(strings.TrimLeft(trim, " \t")))
		} else if trim == "html5" {
			html5()
		} else if trim == "ow:" {
			pyOpenWrite(line)
		} else if trim == "ubb" {
			fmt.Println("#!/usr/bin/env bash")
		} else if trim == "ubp" {
			fmt.Println("#!/usr/bin/env python")
		} else if strings.HasPrefix(trim, "t.t") {
			fmt.Println("func TestFoo (t *testing.T){")
		} else {
			fmt.Println(line)
		}
	}
}

func printNotNilFatal() {
	// print error block
	fmt.Println(`if err != nil{
	log.Fatalf("Failed to do something: %s\n", err)
}`)
}

func printNotNilLog() {
	// print error block
	fmt.Println("if err != nil{")
	fmt.Println(`log.Printf("Failed to do something: %s\n", err)`)
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
	parts := strings.Split(line, " ")
	name := "index"
	if len(parts) > 1 {
		name = parts[1]
	}
	fmt.Printf(`func %s(w http.ResponseWriter, r *http.Request){
	}`, name)
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

func openFile() {
	fmt.Println(`f, err := os.Open(filename)
	if err != nil {
		log.Printf("Unable to open %q: %s\n", filename, err)
	}
	defer f.Close()`)
}

func getURL() {
	fmt.Println(`resp, err := http.Get(link)
	if err != nil {
		log.Printf("Unable to fetch %q: %s\n", link, err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read response from  %q: %s\n", link, err)
		return
	}
	`)
}

func reqStdin() {
	fmt.Println(`stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if stat.Mode()&os.ModeCharDevice != 0 {
		log.Fatal("please pipe in some data")
	}`)
}

func tempFile() {
	fmt.Println(`t, err := ioutil.TempFile("", "temp")
if err != nil{
	log.Fatalf("Unable to create temp file: %s\n", err)
}
fmt.Printf("Created temp file %q\n", t.Name())
defer t.Close()
`)
}
