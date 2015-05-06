package diags

import (
	"fmt"
	"regexp"
)

type ReList []*regexp.Regexp

func (rl ReList) FindStringIndex(s string) []int {
	for _, re := range rl {
		if idx := re.FindStringIndex(s); idx != nil {
			//log.Printf("%s - %q", re, s)
			return idx
		}
	}
	return nil
}

func (rl ReList) MatchString(s string) bool {
	return rl.FindStringIndex(s) != nil
}

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
	ReCKeywords ReList

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
		`\+`, `\*`, `/`, `%`, `\^`, `\?`, `:`,
	}
	reCMultiOps = []string{
		`(?:[^t\n\f\r (]-|-[^\t\n\f\r 0-9a-zA-Z=])`,                                  // -
		`(?:[^\t\n\f\r =!><+-\\*/&|^%]=|=[^\t\n\f\r =])`,                             // =
		`(?:[^\t\n\f\r ><][=!><+-\\*/&|^%]=)`,                                        // *=
		`(?:[^\t\n\f\r |]\||\|[^\t\n\f\r |=])`, `(?:[^\t\n\f\r &]&|&[^\t\n\f\r &=])`, // |, &
		`(?:\S[^& ]&|&[^& ]\S)`, `(?:\S[^| ]\||\|[^| ]\S)`, // &&, ||
		`(?:[^\t\n\f\r <]<|<[^\t\n\f\r <=])`, `(?:[^\t\n\f\r >]>|>[^\t\n\f\r >=])`, // <, <<, >, >>
	}
	ReCOperators ReList
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
