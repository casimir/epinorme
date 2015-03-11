package main

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testCtxt = ErrorContext{}

func TestDiagIdentifier(t *testing.T) {
	Convey("It should detect invalid case end ending number", t, func() {
		So(CheckIdentifier(testCtxt, "someid"), ShouldBeEmpty)
		So(CheckIdentifier(testCtxt, "some_id"), ShouldBeEmpty)

		So(CheckIdentifier(testCtxt, "someId"), ShouldNotBeEmpty)
		So(CheckIdentifier(testCtxt, "Someid"), ShouldNotBeEmpty)
		So(CheckIdentifier(testCtxt, "SomeId"), ShouldNotBeEmpty)

		list := CheckIdentifier(testCtxt, "SomeId")
		So(len(list), ShouldEqual, 1)
		So(list[0].Type, ShouldEqual, ErrCamelCase)

		list2 := CheckIdentifier(testCtxt, "someid2")
		So(len(list2), ShouldEqual, 1)
		So(list2[0].Type, ShouldEqual, WarnIdentifier)
	})
}

func TestDiagLine(t *testing.T) {
	Convey("It should detect lines too long", t, func() {
		So(CheckLine(testCtxt, l("Some line")), ShouldBeEmpty)

		list := CheckLine(testCtxt, l(strings.Repeat("x", 81)))
		So(len(list), ShouldEqual, 1)
		So(list[0].Type, ShouldEqual, ErrTooMuchColumn)
	})

	Convey("It should detect wrong indentation", t, func() {
		So(CheckLine(testCtxt, l("Some line")), ShouldBeEmpty)
		So(CheckLine(testCtxt, l("  Some line")), ShouldBeEmpty)
		So(CheckLine(testCtxt, l("    Some line")), ShouldBeEmpty)

		list := CheckLine(testCtxt, l("   heading spaces"))
		So(len(list), ShouldEqual, 1)
		So(list[0].Type, ShouldEqual, WarnBadIndent)
	})

	Convey("It should detect EOL withespaces", t, func() {
		So(CheckLine(testCtxt, l("Some line")), ShouldBeEmpty)

		listSpace := CheckLine(testCtxt, l("extra spaces   "))
		So(len(listSpace), ShouldEqual, 1)
		So(listSpace[0].Type, ShouldEqual, ErrExtraWS)
		listTab := CheckLine(testCtxt, l("extra tabs		"))
		So(len(listTab), ShouldEqual, 1)
		So(listTab[0].Type, ShouldEqual, ErrExtraWS)
	})

	Convey("It should detect multiple errors", t, func() {
		badLine := "   " + strings.Repeat("x", 80) + "	"
		So(len(CheckLine(testCtxt, l(badLine))), ShouldEqual, 3)
	})
}

func TestDiagFunction(t *testing.T) {
	Convey("It should detect wrong bracket placement", t, func() {
		file, _ := NewFile("_test/func_brackets.c")

		So(CheckFunction(testCtxt, file.Funcs[0]), ShouldBeEmpty)

		elist := CheckFunction(testCtxt, file.Funcs[1])
		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrBracketPlacement)
	})

	Convey("It should detect extra arguments", t, func() {
		file, _ := NewFile("_test/func_args.c")

		So(CheckFunction(testCtxt, file.Funcs[0]), ShouldBeEmpty)

		elist := CheckFunction(testCtxt, file.Funcs[1])
		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrTooMuchArg)
	})

	Convey("It should detect extra lines", t, func() {
		file, _ := NewFile("_test/func_lines.c")

		So(CheckFunction(testCtxt, file.Funcs[0]), ShouldBeEmpty)

		elist := CheckFunction(testCtxt, file.Funcs[1])
		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrTooMuchLine)
		So(elist[0].Line, ShouldEqual, 57)
	})
}

func TestDiagHeader(t *testing.T) {
	Convey("It should detect header errors", t, func() {
		goodHeader := []string{
			"/*",
			"** µFILENAMEµ for µPROJECTNAMEµ in µPATHFILEµ",
			"** ",
			"** Made by µNAMEµ",
			"** Login   <µLOGINµ@epitech.eu>",
			"** ",
			"** Started on  µCREATDAYµ µNAMEµ",
			"** Last update µLASTUPDATEµ µLOGINLASTµ",
			"*/",
			"",
		}
		So(CheckHeader(testCtxt, goodHeader), ShouldBeEmpty)

		badHeader1 := []string{"/*", "** ", "*/"}
		So(CheckHeader(testCtxt, badHeader1), ShouldNotBeEmpty)
		badHeader2 := []string{"/*", "**", "** ", "** ", "** ", "** ", "** ", "** ", "*/"}
		So(CheckHeader(testCtxt, badHeader2), ShouldNotBeEmpty)
		badHeader3 := []string{"/** ", "** ", "** ", "** ", "** ", "** ", "** ", "** ", "*/"}
		So(CheckHeader(testCtxt, badHeader3), ShouldNotBeEmpty)
		badHeader4 := []string{"/*", "** ", "** ", "** ", "** ", "** ", "** ", "** ", "**/"}
		So(CheckHeader(testCtxt, badHeader4), ShouldNotBeEmpty)
	})
}

func TestDiagFile(t *testing.T) {
	Convey("It should detect extra functions", t, func() {
		file, _ := NewFile("_test/file_6funcs.c")

		elist := CheckFile(testCtxt, file)
		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrTooMuchFunc)
	})

	Convey("It should handle unknown file types", t, func() {
		file, _ := NewFile("_test/file.notc")
		elist := CheckFile(testCtxt, file)

		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrUnkownFileType)
	})
}

func l(s string) Line {
	return Line{42, s}
}
