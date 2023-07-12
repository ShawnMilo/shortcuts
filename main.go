package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

// Don't overwrite *this* file while editing.
var exemptKey = "1jT18gRquHb4UJXk6XG169YZJ10"
var fileType string

var replace = map[string]func(){
	"_adoca":    adocArticle,
	"_adocb":    adocBook,
	"dbg":       dbg,
	"flagsh":    flagsh,
	"html5":     html5,
	"ubb":       bash,
	"ubp":       python,
	"inotify":   inotify,
}

var update = map[string]func(string){
	":ksuid:": getKSUID,
}

var modify = map[string]string{
	"_ctx,":   "ctx context.Context,",
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

		var indentation string
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if indent > 0 {
			indentation = line[:indent]
		}
		trim := strings.TrimSpace(line)
		if f, found := replace[trim]; found {
			// don't lose indentation
			fmt.Printf(indentation)
			f()
			continue
		}

		var changed bool
		for pre, f := range update {
			if strings.Contains(trim, pre) {
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
			// Don't re-add indentation if the line is blank.
			fmt.Println(trim)
			continue
		} else {
			lastBlank = false
		}
		fmt.Println(indentation + trim)
	}
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
	fmt.Printf(`lg.Debugf("#%s%s %%v", x)`+"\n", weekday, day)
}

func dp() {
	weekday := string([]rune(time.Now().Weekday().String())[0])
	day := time.Now().Format("02")
	fmt.Printf(`print(f'#%s%s: ')`+"\n", weekday, day)
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
v0.1.0, %s
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

func getKSUID(l string) {
	fmt.Println(strings.Replace(l, ":ksuid:", ksuid.New().String(), 1))
}

func inotify() {
	fmt.Println(`#!/usr/bin/env bash

while true
do
    export filename=$(inotifywait -t 3600 -r -e close_write --format %w%f --include '\.py$' .)
    if [ "$filename" == '' ]; then
        break
    fi
    ./$filename
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
