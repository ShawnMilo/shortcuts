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
