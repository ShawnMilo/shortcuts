package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var replace = map[string]func(){
	"nnf":       nnf,
	"nnl":       nnl,
	"openFile":  openFile,
	"readFile":  readFile,
	"getURL":    getURL,
	"reqStdin":  reqStdin,
	"goMain":    goMain,
	"tempFile":  tempFile,
	"serveHTTP": serveHTTP,
	"pymain":    pyMain,
	"html5":     html5,
	"now":       now,
	"ubb":       bash,
	"ubp":       python,
	"gomain":    goMain,
	"flagsh":    flagsh,
	"dummyType": dummyType,
}

var update = map[string]func(string){
	"fpl(":  fpl,
	"lpf(":  lpf,
	"lpl(":  lpl,
	"fpf(":  fpf,
	"hfunc": hfunc,
	"ow:":   pyOpenWrite,
	"ul":    ul,
}

func main() {
	if os.Getenv("shortcuts") == "off" {
		return
	}

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

		if f, found := replace[trim]; found {
			f()
			continue
		}

		var replaced bool
	DONE:
		for pre := range update {
			if strings.HasPrefix(trim, pre) {
				update[pre](trim)
				replaced = true
				break DONE
			}

		}
		if !replaced {
			fmt.Println(line)
		}
	}

}

func nnf() {
	fmt.Println(`if err != nil{
	log.Fatalf("Failed to do something: %s\n", err)
	}`)
}

func nnl() {
	fmt.Println(`if err != nil{
	log.Printf("Failed to do something: %s\n", err)
	}`)
}

func lpf(line string) {
	fmt.Println(strings.Replace(line, "lpf(", "log.Printf(", 1))
}

func fpf(line string) {
	fmt.Println(strings.Replace(line, "fpf(", "fmt.Printf(", 1))
}

func lpl(line string) {
	fmt.Println(strings.Replace(line, "lpl(", "log.Println(", 1))
}

func fpl(line string) {
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

func serveHTTP() {
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
        <link rel="stylesheet" href="https://aoeus.com/milligram.min.css" type="text/css">
		<meta name="viewport" content="width-device-width, initial-scale=1">
        <style type="text/css">
            body{
                margin: 40px auto;
                max-width: 650px;
                line-height: 1.6;
                font-size: 18px; 
                color: #444;
                padding:0 10px;
            }   
            h1,h2,h3{
                line-height:1.2
            }   
        </style>
	</head>
	<body>
        <article>
            <h1>topic</h1>
            <p>content</p>
        </article>
	</body>
</html>
`)
}

func hfunc(line string) {
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

func ul(line string) {
	trim := strings.TrimSpace(line)
	margin := len(line) - len(strings.TrimLeft(trim, " \t"))
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

func readFile() {
	fmt.Println(`b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to open %q: %s\n", filename, err)
	}`)
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

func now() {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
}

func bash() {
	fmt.Println("#!/usr/bin/env bash")
}

func python() {
	fmt.Println("#!/usr/bin/env python")
}

func flagsh() {

	fmt.Println(`#!/usr/bin/env bash

flag=$(mktemp)
touch $flag

while true; do
sleep 5
    find . -mmin -1 -name '*.go' 2>>/dev/null | while read file; do
        if [[ "$file" -nt $flag ]]; then
            if [[ "$file" == "$flag" ]]; then
                continue
            fi
            echo "$file was updated"
            touch $flag
        fi
    done
done
`)
}

func dummyType() {
	fmt.Println(`type dummy struct {
    thing string
    size int
    color string
}
`)
}
