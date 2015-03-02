package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s FILES...\n", os.Args[0])
		os.Exit(1)
	}
	for _, it := range flag.Args() {
		file, err := NewFile(it)
		if err != nil {
			log.Print(err)
			continue
		}
		ctxt := ErrorContext{File: file.Name}
		for _, e := range CheckFile(ctxt, file) {
			fmt.Println(e)
		}
	}
}
