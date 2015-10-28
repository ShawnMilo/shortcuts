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
		} else if strings.HasPrefix(line, "lpl(") {
			lPrintln(t)
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
