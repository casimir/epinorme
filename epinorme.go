package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/casimir/epinorme/diags"
)

var (
	aMark    = flag.Bool("mark", false, "Use mark mode")
	aNoErr   = flag.Bool("noerr", false, "Hide errors")
	aNoWarn  = flag.Bool("nowarn", false, "Hide warnings")
	aProject = flag.Bool("project", false, "Use project mode")
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] FILES...\n", os.Args[0])
		os.Exit(1)
	}
	diags.PrintErr = !*aNoErr
	diags.PrintWarn = !*aNoWarn
	if *aProject {
		runProjectMode()
	} else {
		runFileMode()
	}
}

func runFileMode() {
	for _, it := range flag.Args() {
		file, err := diags.NewFile(it)
		if err != nil {
			log.Print(err)
			continue
		}
		ctxt := diags.ErrorContext{File: file.Name}
		for _, e := range diags.CheckFile(ctxt, file) {
			if e.ShouldPrint() {
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
			p := diags.NewProject(proj)
			if *aMark {
				fmt.Printf("%s: %d\n", p.Name, p.Note())
			} else {
				fmt.Println(p)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
