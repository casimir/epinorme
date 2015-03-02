package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorContext(t *testing.T) {
	testCtxt := ErrorContext{
		File:   "file.c",
		Line:   42,
		Column: 12,
	}
	errStr := "file.c:42:12:"
	warnStr := "file.c:42: warning: "

	Convey("It should generate a valid error", t, func() {
		err1 := testCtxt.NewError(ErrBracketPlacement)
		expected1 := errStr + ErrMessages[ErrBracketPlacement]
		So(err1.Error(), ShouldEqual, expected1)

		err2 := testCtxt.NewError(WarnIdentifier)
		expected2 := warnStr + ErrMessages[WarnIdentifier]
		So(err2.Error(), ShouldEqual, expected2)

		err3 := testCtxt.NewError(ErrType(-1))
		expected3 := errStr + ErrMessages[ErrUnknown]
		So(err3.Error(), ShouldEqual, expected3)
	})
}
