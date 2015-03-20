package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	aErr     = flag.Bool("e", true, "Show errors")
	aMark    = flag.Bool("mark", false, "Use mark mode")
	aProject = flag.Bool("project", false, "Use project mode")
	aStats   = flag.Bool("stats", false, "Print some file statistics")
	aWarn    = flag.Bool("w", true, "Show warnings")
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
				if e.Type > warnBegin && *aWarn {
					fmt.Println(e)
				} else if e.Type < warnBegin && *aErr {
					fmt.Println(e)
				}
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
