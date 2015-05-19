package diags

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
		So(CheckLine(testCtxt, l("Some line"), false), ShouldBeEmpty)

		list := CheckLine(testCtxt, l(strings.Repeat("x", 81)), false)
		So(len(list), ShouldEqual, 1)
		So(list[0].Type, ShouldEqual, ErrTooMuchColumn)
	})

	Convey("It should detect wrong indentation", t, func() {
		So(CheckLine(testCtxt, l("Some line"), false), ShouldBeEmpty)
		So(CheckLine(testCtxt, l("  Some line"), false), ShouldBeEmpty)
		So(CheckLine(testCtxt, l("    Some line"), false), ShouldBeEmpty)

		list := CheckLine(testCtxt, l("   heading spaces"), false)
		So(len(list), ShouldEqual, 1)
		So(list[0].Type, ShouldEqual, WarnBadIndent)
	})

	Convey("It should detect EOL withespaces", t, func() {
		So(CheckLine(testCtxt, l("Some line"), false), ShouldBeEmpty)

		listSpace := CheckLine(testCtxt, l("extra spaces   "), false)
		So(len(listSpace), ShouldEqual, 1)
		So(listSpace[0].Type, ShouldEqual, ErrExtraWS)
		listTab := CheckLine(testCtxt, l("extra tabs		"), false)
		So(len(listTab), ShouldEqual, 1)
		So(listTab[0].Type, ShouldEqual, ErrExtraWS)
	})

	Convey("It should detect missing whitespaces (for keywords and operators)", t, func() {
		good := []string{
			"do_something()",
			"do",
			"undo()",
			"return EXIT_FAILURE;",
			"while (42)",
			"a + b",
			"f(\"+ab\")",
			"if (test == 1)",
			"val *= 42)",
			"val <<= 4)",
			"true || false",
			"return (-42);",
			"i++;",
			"!predicat",
		}
		bad := []string{
			"return(1);",
			"while(42)",
			"else if(-42)",
			"a* b",
			"if (test==1)",
			"val*= 42)",
			"val <<=4)",
			"false&& true",
			"12-42",
			"! nope",
		}

		for _, it := range good {
			el := CheckLine(testCtxt, l(it), true)
			So(el, ShouldBeEmpty)
		}
		for _, it := range bad {
			el := CheckLine(testCtxt, l(it), true)
			So(len(el), ShouldEqual, 1)
			So(el[0].Type, ShouldEqual, ErrMissingSpace)
		}
	})

	Convey("It should detect wrong ponctuation placement", t, func() {
		lines := []string{
			"toto , titi",
			"toto,titi",
			"toto;titi",
		}

		for _, it := range lines {
			el := CheckLine(testCtxt, l(it), false)
			So(len(el), ShouldEqual, 1)
			So(el[0].Type, ShouldEqual, ErrPonctPlacement)
			So(el[0].Column, ShouldEqual, 5)
		}

		el4 := CheckLine(testCtxt, l("printf(\";\");"), false)
		So(el4, ShouldBeEmpty)
	})

	Convey("It should detect multiple errors", t, func() {
		badLine := "   " + strings.Repeat("x", 80) + "	"
		So(len(CheckLine(testCtxt, l(badLine), false)), ShouldEqual, 3)
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

		So(len(file.Funcs), ShouldEqual, 2)
		So(CheckFunction(testCtxt, file.Funcs[0]), ShouldBeEmpty)

		elist := CheckFunction(testCtxt, file.Funcs[1])
		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrTooMuchArg)
		So(elist[0].Line, ShouldEqual, 6)
		So(elist[0].Column, ShouldEqual, 15)
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
		So(elist[0].Line, ShouldEqual, 31)
	})

	Convey("It should handle unknown file types", t, func() {
		file, _ := NewFile("_test/file.notc")
		elist := CheckFile(testCtxt, file)

		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, WarnUnknownFileType)
	})

	Convey("It should report header instruction in source files", t, func() {
		file, _ := NewFile("_test/proto.c")
		elist := CheckFile(testCtxt, file)

		So(len(elist), ShouldEqual, 1)
		So(elist[0].Type, ShouldEqual, ErrHeaderInSource)
	})
}

func l(s string) Line {
	return Line{42, s}
}
