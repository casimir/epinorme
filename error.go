package main

import "fmt"

type ErrType int

// TODO handle comments
const (
	ErrUnknown = iota

	ErrBracketPlacement
	ErrCamelCase
	ErrExtraWS
	ErrMissingBlankLine // TODO
	ErrMissingSpace     // TODO
	ErrTooMuchBlankLine // TODO
	ErrTooMuchInstruc   // TODO

	errBeginSerious
	ErrMissingHeader
	ErrTooMuchArg
	ErrTooMuchColumn
	ErrTooMuchFunc
	ErrTooMuchLine

	warnBegin
	WarnBadIndent
	WarnIdentifier
	WarnSyscallRetCheck // TODO
	WarnUnknownFileType
)

var ErrMessages = map[ErrType]string{
	ErrUnknown:          "Unknown error",
	ErrBracketPlacement: "Wrong bracket placement",
	ErrCamelCase:        "Bad identifier casing",
	ErrExtraWS:          "Extra whitespace(s) at EOL",
	ErrMissingHeader:    "Missing header",
	ErrMissingSpace:     "Missing space after keyword",
	ErrTooMuchArg:       "More than 4 arguments",
	ErrTooMuchColumn:    "More than 80 columns",
	ErrTooMuchFunc:      "More than 5 functions",
	ErrTooMuchLine:      "More than 25 lines",
	WarnUnknownFileType: "Unknown file type",

	WarnBadIndent:  "Possibly wrong indentation",
	WarnIdentifier: "Poorly named identifier",
}

type Error struct {
	Type   ErrType
	File   string
	Line   int
	Column int
	What   string
}

func (e Error) String() string {
	if e.Type > warnBegin {
		return fmt.Sprintf("%s:%d: warning: %s", e.File, e.Line, e.What)
	}
	return fmt.Sprintf("%s:%d:%d:%s", e.File, e.Line, e.Column, e.What)
}

func (e Error) Error() string {
	return e.String()
}

type ErrorContext struct {
	File   string
	Line   int
	Column int
}

func (c ErrorContext) NewError(et ErrType) Error {
	if _, ok := ErrMessages[et]; !ok {
		et = ErrUnknown
	}
	return Error{
		Type:   et,
		File:   c.File,
		Line:   c.Line,
		Column: c.Column,
		What:   ErrMessages[et],
	}
}
