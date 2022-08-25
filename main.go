package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Don't overwrite *this* file while editing.
var exemptKey = "1jT18gRquHb4UJXk6XG169YZJ10"
var fileType string

var replace = map[string]func(){
	"_adoca":    adocArticle,
	"_adocb":    adocBook,
	"_audit":    audit,
	"bootstrap": bootstrap,
	"dbg":       dbg,
	"flagsh":    flagsh,
	"getURL":    getURL,
	"gomain":    goMain,
	"html5":     html5,
	"jm":        jsonMarshal,
	"ju":        jsonUnmarshal,
	"_lg":       lg,
	"nnf":       nnf,
	"nnl":       nnl,
	"openFile":  openFile,
	"pymain":    pyMain,
	"rb":        readBody,
	"readFile":  readFile,
	"reqStdin":  reqStdin,
	"serveHTTP": serveHTTP,
	"src":       src,
	"tempFile":  tempFile,
	"ubb":       bash,
	"ubp":       python,
	"watcher":   watcher,
}

var update = map[string]func(string){
	":cb:":   markdownCheckboxes,
	"fpf(":   fpf,
	"fpl(":   fpl,
	"hfunc":  hfunc,
	":json:": formatJSON,
	"lpf(":   lpf,
	"lpl(":   lpl,
	"==day":  adocDay,
	"===now": adocNow,
	"ow:":    pyOpenWrite,
	":tb:":   markdownTable,
	"ul":     ul,
}

var modify = map[string]string{
	"_ctx,":   "ctx context.Context,",
	"_ctc":    "Cracking the Cryptic",
	":sg:":    "😎",
	":un:":    "😒",
	":check:": "✅",
	":x:":     "❌",
	":cc:":    "☑",
	":d:":     func() string { return string(time.Now().Weekday().String()[0]) + time.Now().Format("02") }(),
	":ce:":    "☐",
	":cx:":    "☒",
	":boom:":  "💥",
	":cool:":  "🆒",
	":ok:":    "🆗",
	":now:":   func() string { return time.Now().Format("2006-01-02 15:04:05") }(),
	":cr:": func() string {
		now := time.Now().Format(time.RFC3339)
		return fmt.Sprintf("---\ncreated: %s\nmodified: %s\n---\n\n# Title", now, now)
	}(),
	":mod:": func() string { return time.Now().Format(time.RFC3339) }(),
}

func getFileType(line string) string {
	if strings.Contains(line, "env python") {
		return "python"
	}
	if strings.Contains(line, "env bash") {
		return "bash"
	}
	if strings.Contains(line, "env fish") {
		return "fish"
	}
	if strings.Contains(line, "doctype") {
		return "html"
	}
	if strings.HasPrefix(line, "package") {
		return "go"
	}
	if strings.HasPrefix(line, "= ") {
		return "adoc"
	}
	return ""
}

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if stat.Mode()&os.ModeCharDevice != 0 {
		log.Fatal("please pipe in some data")
	}

	s := bufio.NewScanner(os.Stdin)
	var lastBlank bool

	// Do not overwrite files containing this flag.
	// This is mostly to protect *this* file.
	var exempt bool

	var count int
	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, exemptKey) {
			exempt = true
		}

		//  get file type
		if count == 0 {
			if len(os.Args) == 2 {
				fileType = os.Args[1]
			} else {
				fileType = getFileType(line)
			}
		}

		count++

		if exempt {
			fmt.Println(strings.TrimSuffix(line, "\n"))
			continue
		}

		trim := strings.TrimSpace(line)

		if f, found := replace[trim]; found {
			f()
			continue
		}

		var changed bool
		for pre, f := range update {
			if strings.HasPrefix(trim, pre) {
				f(line)
				changed = true
				break
			}

		}
		if changed {
			continue
		}

		for pre, post := range modify {
			if strings.Contains(line, pre) {
				fmt.Println(strings.Replace(line, pre, post, 5))
				changed = true
			}
		}
		if changed {
			continue
		}

		if trim == "" {
			if lastBlank {
				continue
			}
			lastBlank = true
		} else {
			lastBlank = false
		}
		fmt.Println(line)
	}
}

func audit() {
	fmt.Printf(`// Last audit %s by skm.`, time.Now().Format("2006-01-02"))
}

func lg() {
	fmt.Println(`lg := logger.FromContext(ctx)`)
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

func dbg() {
	if fileType == "go" {
		dg()
		return
	}
	if fileType == "py" {
		dp()
		return
	}
}

func dg() {
	weekday := string([]rune(time.Now().Weekday().String())[0])
	day := time.Now().Format("02")
	fmt.Printf(`lg.Debugf("%s%s %%v", x)`+"\n", weekday, day)
}

func dp() {
	weekday := string([]rune(time.Now().Weekday().String())[0])
	day := time.Now().Format("02")
	fmt.Printf(`print(f'%s%s: ')`+"\n", weekday, day)
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
	fmt.Println(`#!/usr/bin/env python3
"""
You should probably write something here.
"""

def main():
    """
    Do the thing.
    """
    print("python")

if __name__ == '__main__':
    main()
`)
}

func bootstrap() {
	fmt.Println(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Bootstrap demo</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-gH2yIJqKdNHPEq0n4Mqa/HGKIhSkIHeL5AyhkYV8i59U5AR6csBvApHHNl/vI1Bx" crossorigin="anonymous">
  </head>
  <body>
    <h1>Hello, world!</h1>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/js/bootstrap.bundle.min.js" integrity="sha384-A3rJD856KowSb7dwlZdYEkO39Gagi7vIsF0jrRAoQmDKKtQBHUuLZ9AsSv4jD4Xa" crossorigin="anonymous"></script>
  </body>
</html>`)
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

func adocArticle() {
	name := fromGitConfig("user.name")
	email := fromGitConfig("user.email")
	date := time.Now().Format("2006-01-02")
	fmt.Printf(`= Article Title
%s <%s>
v0.1, %s
:doctype: article
:source-highlighter: pygments
:toc:
:icons: font`, name, email, date)
}

func adocBook() {
	name := fromGitConfig("user.name")
	email := fromGitConfig("user.email")
	date := time.Now().Format("2006-01-02")
	fmt.Printf(`= Book Title
%s <%s>
v0.1.0, %s
:doctype: book
:source-highlighter: pygments
:toc:
:icons: font

[index]
== Index`, name, email, date)
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

func markdownCheckboxes(line string) {
	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		fmt.Println(line)
		return
	}
	count, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		fmt.Println(line)
		return
	}
	for i := 0; i < count; i++ {
		fmt.Println("- [ ] ")
	}
}

func markdownTable(line string) {
	parts := strings.Split(line, ":")
	if len(parts) < 4 {
		fmt.Println(line)
		return
	}
	rows, err := strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		fmt.Println(line)
		return
	}
	columns, err := strconv.Atoi(strings.TrimSpace(parts[3]))
	if err != nil {
		fmt.Println(line)
		return
	}

	fmt.Printf("|")
	fmt.Println(strings.Repeat(" aoeus |", columns))

	fmt.Printf("|")
	fmt.Println(strings.Repeat(" --- |", columns))

	for i := 0; i < rows; i++ {
		fmt.Printf("|")
		for j := 0; j < columns; j++ {
			fmt.Printf(" aoeus |")
		}
		fmt.Printf("\n")
	}
	fmt.Println()
}

func formatJSON(line string) {
	b := []byte(strings.TrimSpace(line)[6:]) // strip off :json:
	var thing map[string]interface{}
	err := json.Unmarshal(b, &thing)
	if err != nil {
		fmt.Println(line)
		return
	}
	j, err := json.MarshalIndent(thing, "", "    ")
	if err != nil {
		fmt.Println(line)
		return
	}
	fmt.Printf("%s\n", string(j))
}

func ul(line string) {
	if fileType != "html" {
		fmt.Println(line)
		return
	}
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

func adocNow(line string) {
	fmt.Println(strings.Replace(line, "===now", "=== "+time.Now().Format("2006-01-02 15:04:05"), 1))
}

func adocDay(line string) {
	fmt.Println(strings.Replace(line, "==day", "== "+time.Now().Format("Monday, January 2, 2006"), 1))
}

func bash() {
	fmt.Println("#!/usr/bin/env bash")
}

func python() {
	fmt.Println("#!/usr/bin/env python")
}

func src() {
	fmt.Println(`[source,go]
----
include::filename.go[]
----

`)
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

func jsonMarshal() {
	fmt.Println(`b, err := json.Marshal(x)
    if err != nil{
        log.Fatalf("Failed to do something: %s\n", err)
	 } `)
}

func jsonUnmarshal() {
	fmt.Println(`var x thing
err = json.Unmarshal(b, &x)
return x, err
`)
}

func readBody() {
	fmt.Println(`b, err := ioutil.ReadAll(r.Body)
	 if err != nil {
		 return err
	 }
	 defer r.Body.Close()`)
}

func watcher() {
	fmt.Println(`#!/usr/bin/env bash

bin=monkey
flag=$(mktemp)
code=main.go

function cleanup() {
    pkill $bin
    rm -f $bin
    rm -f $flag
    exit
}

trap cleanup SIGINT SIGTERM

while true
do
    if [ ! "$code" -nt $flag ]; then
        sleep 5
        continue
    fi
    go build -o $bin $code || { touch $flag; continue; }
    pkill $bin
    ./$bin &
    sleep 1
    ./tests.sh
    touch $flag
    sleep 15
done`)
}

func fromGitConfig(value string) string {
	// git config --list
	cmd := exec.Command("git", "config", "--list")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, value) {
			return strings.TrimSpace(line[len(value)+1:])
		}
	}
	return ""
}
