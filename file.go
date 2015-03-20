package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	ipath "path"
	"regexp"
	"strings"
)

var (
	reFuncBeg = regexp.MustCompile(`^\w+\s+\**\w+\(`)
	reFuncEnd = regexp.MustCompile(`\)`)
	reFunc    = regexp.MustCompile(`(?ms)^(\w+)\s+(\**)(\w+)\((.*)\)`)
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

func (f File) ListFuncs() []string {
	var ret []string
	for _, it := range f.Funcs {
		ret = append(ret, it.Name)
	}
	return ret
}

func (f File) String() string {
	parts := []string{
		f.Name,
		fmt.Sprintf("→ Functions (%d): %v\n", len(f.Funcs), f.ListFuncs()),
	}
	return strings.Join(parts, "\n")
}

func NewFile(path string) (file File, err error) {
	file = File{
		Name: path,
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

		if reFuncBeg.MatchString(line) {
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
		RetType   string
		Name      string
		Args      []Argument
		Lines     []Line
		protoSize int
	}
)

func (f Function) innerLines() (first, last int) {
	first = f.protoSize
	last = first
	if strings.HasPrefix(f.Lines[first].str, "{") {
		last -= 1
	}
	for i, it := range f.Lines[first:] {
		if strings.HasPrefix(it.str, "}") {
			last += i - 1
			break
		}
	}
	return
}

func newFunction(s *bufio.Scanner, n *int) (fn Function) {
	var accu []string
	for !reFuncEnd.MatchString(s.Text()) {
		accu = append(accu, s.Text())
		if !s.Scan() {
			log.Fatal("Fucked up function")
		}
	}
	accu = append(accu, s.Text())
	parts := reFunc.FindStringSubmatch(strings.Join(accu, "\n"))
	fn = Function{
		RetType:   parts[1] + parts[2],
		Name:      parts[3],
		protoSize: len(accu),
	}
	for _, it := range accu {
		fn.Lines = append(fn.Lines, Line{*n, it})
		*n++
	}

	args := parts[4]
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
	*n--
	return fn
}
