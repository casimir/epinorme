package diags

import "strings"

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
