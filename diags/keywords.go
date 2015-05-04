package diags

import (
	"fmt"
	"regexp"
)

var (
	CKeywords = []string{
		"auto",
		"const",
		"do",
		"double",
		"else",
		"enum",
		"extern",
		"if",
		"register",
		"return",
		"signed",
		"static",
		"struct",
		"typedef",
		"union",
		"unsigned",
		"volatile",
		"while",
	}
	ReCKeywords []*regexp.Regexp

	CKeywordsForbidden = []string{
		"break",
		"continue",
		"case",
		"default",
		"for",
		"goto",
		"switch",
	}

	COperators = []string{
		`\+`, "-", `\*`, "/", "%",
		"==", "!=", ">", "<", ">=", "<=",
		"&", `\|`, `\^`, "<<", ">>",
		"&&", `\|\|`, "!",
		"=", `\+=`, `-\+`, `\*=`, "/=", "%=", "<<=", "<<=", "&=", `\^=`, "!=",
		`\?`, ":",
	}
	ReCOperators []*regexp.Regexp
)

func init() {
	for _, it := range CKeywords {
		re := fmt.Sprintf(`\b(%s)\b\S`, it)
		ReCKeywords = append(ReCKeywords, regexp.MustCompile(re))
	}

	for _, it := range COperators {
		re := fmt.Sprintf(`(\S%s|%s\S)`, it, it)
		ReCOperators = append(ReCOperators, regexp.MustCompile(re))
	}
}
