package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewFile(t *testing.T) {
	Convey("It should parse a source file and initialize the struct", t, func() {
		f, err := NewFile("_test/func_lines.c")
		So(err, ShouldBeNil)
		So(f.Type, ShouldEqual, CSource)
		So(f.Name, ShouldEqual, "func_lines.c")
		So(len(f.Funcs), ShouldEqual, 2)

		fn0 := f.Funcs[0]
		f0, l0 := f.Funcs[0].innerLines()
		So(fn0.Name, ShouldEqual, "good")
		So(fn0.Lines[0].n, ShouldEqual, 1)
		So(fn0.Lines[f0].n, ShouldEqual, 3)
		So(fn0.Lines[l0].n, ShouldEqual, 27)

		fn1 := f.Funcs[1]
		f1, l1 := f.Funcs[1].innerLines()
		So(fn1.Name, ShouldEqual, "bad")
		So(fn1.Lines[0].n, ShouldEqual, 30)
		So(fn1.Lines[f1].n, ShouldEqual, 32)
		So(fn1.Lines[l1].n, ShouldEqual, 61)
	})

	Convey("It should handle unkwon file types correctly", t, func() {
		f, err := NewFile("_test/file.notc")
		So(err, ShouldNotBeNil)
		So(f.Type, ShouldEqual, UnkownType)
	})
}
