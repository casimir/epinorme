package diags

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
	// Remove all cstrings to ease further matching.
	line.str = reCString.ReplaceAllString(line.str, "")
	if reExtraWS.MatchString(line.str) {
		ctxt.Column = len(line.str)
		ret = append(ret, ctxt.NewError(ErrExtraWS))
	}
	for _, it := range ReCKeywords {
		if loc := it.FindStringIndex(line.str); loc != nil {
			ctxt.Column = loc[0]
			ret = append(ret, ctxt.NewError(ErrMissingSpace))
		}
	}
	if idxs := reBadPonct.FindStringIndex(line.str); len(idxs) > 0 {
		ctxt.Column = idxs[0] + 1
		ret = append(ret, ctxt.NewError(ErrPonctPlacement))
	}
	return ret
}
