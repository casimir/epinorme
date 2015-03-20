package main

import (
	"fmt"
	"log"
	"os"
	ipath "path"
	"path/filepath"
	"strings"
)

type Project struct {
	Name   string
	Path   string
	Files  []string
	Errors []Error
}

func (p Project) Note() int {
	ret := 0
	counted := map[ErrType]bool{}
	for _, it := range p.Errors {
		if it.Type > ErrUnknown && it.Type < errBeginSerious && !counted[it.Type] {
			ret -= 1
			counted[it.Type] = true
		} else if it.Type > errBeginSerious && it.Type < warnBegin && !counted[it.Type] {
			ret -= 5
			counted[it.Type] = true
		}
	}
	return ret
}

func (p Project) String() string {
	if len(p.Name) == 0 {
		p.Name = "Note"
	}

	if *aMark {
		return fmt.Sprintf("%s: %d", p.Name, p.Note())
	}

	var strs []string
	for _, it := range p.Errors {
		strs = append(strs, it.String())
	}
	return strings.Join(strs, "\n")
}

func NewProject(path string) Project {
	p := Project{
		Name: ipath.Base(path),
		Path: path,
	}
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Failed to access file: %s", path)
		}
		if info == nil || info.IsDir() {
			return nil
		}

		p.Files = append(p.Files, path)
		ctxt := ErrorContext{File: path}
		f, ferr := NewFile(path)
		if ferr != nil {
			return err
		}
		p.Errors = append(p.Errors, CheckFile(ctxt, f)...)
		return nil
	}
	filepath.Walk(path, fn)
	return p
}
