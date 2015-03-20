package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	aNote    = flag.Bool("note", false, "Use note mode")
	aProject = flag.Bool("project", false, "Use project mode")
	aStats   = flag.Bool("stats", false, "Print some file statistics")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] FILES...\n", os.Args[0])
		os.Exit(1)
	}
	if *aProject {
		runProjectMode()
	} else {
		runFileMode()
	}
}

func runFileMode() {
	for _, it := range flag.Args() {
		file, err := NewFile(it)
		if err != nil {
			log.Print(err)
			continue
		}
		if *aStats {
			fmt.Print(file)
		} else {
			ctxt := ErrorContext{File: file.Name}
			for _, e := range CheckFile(ctxt, file) {
				fmt.Println(e)
			}
		}
	}
}

func runProjectMode() {
	var wg sync.WaitGroup
	wg.Add(flag.NArg())

	for _, it := range flag.Args() {
		proj := it
		go func() {
			fmt.Println(NewProject(proj))
			wg.Done()
		}()
	}
	wg.Wait()
}
