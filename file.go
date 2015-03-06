package main

import (
	"bufio"
	"fmt"
	"os"
	ipath "path"
	"regexp"
	"strings"
)

var (
	reFunc    = regexp.MustCompile(`^(\w+)\s+(\w+)\((.*)\)`)
	reInclude = regexp.MustCompile(`^#\s*include\s+(.*)`)
)

type FileType int

const (
	UnkownType = FileType(iota)
	CHeader
	CSource
)

func getFileType(ext string) FileType {
	switch ext {
	case ".h":
		return CHeader
	case ".c":
		return CSource
	}
	return UnkownType
}

type (
	Line struct {
		n   int
		str string
	}

	File struct {
		Name     string
		Type     FileType
		Header   []string
		Includes []Line
		Funcs    []Function
		Fillers  []Line
	}
)

func NewFile(path string) (file File, err error) {
	file = File{
		Name: ipath.Base(path),
		Type: getFileType(ipath.Ext(path)),
	}
	if file.Type == UnkownType {
		ext := ipath.Ext(path)
		if len(ext) > 1 {
			ext = ext[1:]
		}
		err = fmt.Errorf("Unknown extension (%s): %s", ext, path)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for i := 1; s.Scan(); i++ {
		line := s.Text()
		if i <= 10 {
			file.Header = append(file.Header, line)
		}

		if reFunc.MatchString(line) {
			file.Funcs = append(file.Funcs, newFunction(s, &i))
			continue
		}
		if reInclude.MatchString(line) {
			file.Includes = append(file.Includes, Line{i, line})
			continue
		}

		file.Fillers = append(file.Fillers, Line{i, line})
	}
	err = s.Err()
	return
}

type (
	Argument struct {
		Type string
		Name string
	}

	Function struct {
		RetType string
		Name    string
		Args    []Argument
		Lines   []Line
	}
)

func (f Function) innerLines() (first, last int) {
	for i, it := range f.Lines {
		line := strings.TrimSpace(it.str)
		if strings.HasSuffix(line, "{") {
			first = i + 1
			break
		}
	}
	for i, it := range f.Lines[first:] {
		line := strings.TrimSpace(it.str)
		if line == "}" {
			last = i + first - 1
			break
		}
	}
	return
}

func newFunction(s *bufio.Scanner, n *int) (fn Function) {
	parts := reFunc.FindStringSubmatch(s.Text())
	fn = Function{
		RetType: parts[1],
		Name:    parts[2],
		Lines:   []Line{{*n, s.Text()}},
	}

	// FIXME handle multi-line declaration
	args := parts[3]
	if len(args) > 0 {
		for _, it := range strings.Split(args, ",") {
			arg := Argument{}
			if idx := strings.LastIndex(it, " "); idx != -1 {
				arg.Type = it[:idx]
				arg.Name = it[idx:]
			} else {
				arg.Type = it
			}
			fn.Args = append(fn.Args, arg)
		}
	}
	*n++

	still := true
	for ; s.Scan() && still; *n++ {
		line := s.Text()
		fn.Lines = append(fn.Lines, Line{*n, line})
		still = !strings.HasPrefix(line, "}")
	}
	return fn
}
