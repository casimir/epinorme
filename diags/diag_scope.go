package diags

import (
	"regexp"
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
	first, last := fn.innerLines()
	if last-first > 25 {
		ctxt.Line = fn.Lines[first+25].n
		ret = append(ret, ctxt.NewError(ErrTooMuchLine))
	}
	ret = append(ret, CheckLine(ctxt, fn.Lines[0], false)...)
	for _, it := range fn.Lines[first:] {
		ret = append(ret, CheckLine(ctxt, it, true)...)
	}
	return ret
}
