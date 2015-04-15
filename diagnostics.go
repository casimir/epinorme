package main

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	reBadID     = regexp.MustCompile(`\d+$`)
	reBadPonct  = regexp.MustCompile(`\s[,;]|[,;]\S`)
	reExtraWS   = regexp.MustCompile(`\s+$`)
	reLeadingWS = regexp.MustCompile(`^(\s*)`)
)

func CheckIdentifier(ctxt ErrorContext, id string) []Error {
	ret := []Error{}
	for _, it := range id {
		if unicode.IsUpper(it) {
			ret = append(ret, ctxt.NewError(ErrCamelCase))
			break
		}
	}
	if reBadID.Match([]byte(id)) {
		ret = append(ret, ctxt.NewError(WarnIdentifier))
	}
	return ret
}

func CheckLine(ctxt ErrorContext, line Line) []Error {
	// TODO add column info
	ctxt.Line = line.n
	ret := []Error{}
	if len(line.str) > 80 {
		ctxt.Column = len(line.str)
		ret = append(ret, ctxt.NewError(ErrTooMuchColumn))
	}
	lws := reLeadingWS.FindString(line.str)
	if len(lws)%2 != 0 {
		ctxt.Column = len(lws)
		ret = append(ret, ctxt.NewError(WarnBadIndent))
	}
	if reExtraWS.MatchString(line.str) {
		ctxt.Column = len(line.str)
		ret = append(ret, ctxt.NewError(ErrExtraWS))
	}
	if idxs := reBadPonct.FindStringIndex(line.str); len(idxs) > 0 {
		ctxt.Column = idxs[0] + 1
		ret = append(ret, ctxt.NewError(ErrPonctPlacement))
	}
	return ret
}

func CheckFunction(ctxt ErrorContext, fn Function) []Error {
	// XXX context updating doesn't seem right at all
	ctxt.Line = fn.Lines[0].n
	ret := []Error{}
	ret = append(ret, CheckIdentifier(ctxt, fn.Name)...)
	if fn.Lines[fn.protoSize].str[0] != '{' {
		ret = append(ret, ctxt.NewError(ErrBracketPlacement))
	}
	if len(fn.Args) > 4 {
		ctxt.Line = 0
		arg := fn.Args[4]
		for _, it := range fn.Lines[:fn.protoSize] {
			re := regexp.MustCompile(arg.Type + `\s+\**` + arg.Name)
			if idxs := re.FindStringIndex(it.str); len(idxs) > 0 {
				ctxt.Line = it.n
				ctxt.Column = idxs[0] + 2 // match 1 char before + col [1, ...[
				break
			}
		}
		ret = append(ret, ctxt.NewError(ErrTooMuchArg))
		ctxt.Column = 0
	}
	for _, it := range fn.Args {
		ret = append(ret, CheckIdentifier(ctxt, it.Name)...)
	}
	if first, last := fn.innerLines(); last-first > 25 {
		ctxt.Line = fn.Lines[first+25].n
		ret = append(ret, ctxt.NewError(ErrTooMuchLine))
	}
	for _, it := range fn.Lines {
		ret = append(ret, CheckLine(ctxt, it)...)
	}
	return ret
}

func CheckHeader(ctxt ErrorContext, lines []string) []Error {
	// FIXME specific to C
	// TODO email/login format
	noHeaderErr := []Error{ctxt.NewError(ErrMissingHeader)}
	if len(lines) != 10 || lines[0] != "/*" || lines[8] != "*/" || lines[9] != "" {
		return noHeaderErr
	}
	for i := 1; i < 8; i++ {
		if !strings.HasPrefix(lines[i], "**") {
			return noHeaderErr
		}
	}
	ret := []Error{}
	return ret
}

func CheckFile(ctxt ErrorContext, f File) []Error {
	if f.Type != CSource {
		return []Error{ctxt.NewError(WarnUnknownFileType)}
	}
	ret := []Error{}
	ret = append(ret, CheckIdentifier(ctxt, f.Name)...)
	ret = append(ret, CheckHeader(ctxt, f.Header)...)
	for _, it := range f.Protos {
		ctxt.Line = it.n
		ret = append(ret, ctxt.NewError(ErrHeaderInSource))
	}
	if len(f.Funcs) > 5 {
		ctxt.Line = f.Funcs[5].Lines[0].n
		ret = append(ret, ctxt.NewError(ErrTooMuchFunc))
	}
	for _, it := range f.Funcs {
		ret = append(ret, CheckFunction(ctxt, it)...)
	}
	return ret
}
