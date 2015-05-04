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

	reCSimpleOps = []string{
		`\+`, `-`, `\*`, `/`, `%`, `\^`, `\?`, `:`,
	}
	reCMultiOps = []string{
		`(?:\S[^=!><+-\\*/&|^ ]=|=[^\t\n\f\r =])`,                                  // =
		`(?:[^\t\n\f\r ><][=!><+-\\*/&|^]=|[=!><+-\\*/&|^]=\S)`,                    // .=
		`(?:[^\t\n\f\r |]\||\|[^\t\n\f\r |])`, `(?:[^\t\n\f\r &]&|&[^\t\n\f\r &])`, // |, &
		`(?:\S[^& ]&|&[^& ]\S)`, `(?:\S[^| ]\||\|[^| ]\S)`, // &&, ||
		`(?:\S[^< ]<|<[^< ]\S)`, `(?:\S[^> ]>|>[^> ]\S)`, // <<, >>, <<=, >>=
	}
	// TODO generate UT
	ReCOperators []*regexp.Regexp
)

func init() {
	for _, it := range CKeywords {
		re := fmt.Sprintf(`\b(?:%s)\b\S`, it)
		ReCKeywords = append(ReCKeywords, regexp.MustCompile(re))
	}

	for _, it := range reCSimpleOps {
		re := fmt.Sprintf(`(?:\S%s|%s[^\t\n\f\r =])`, it, it)
		ReCOperators = append(ReCOperators, regexp.MustCompile(re))
	}
	for _, it := range reCMultiOps {
		ReCOperators = append(ReCOperators, regexp.MustCompile(it))
	}
}
