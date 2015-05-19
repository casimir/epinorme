package diags

import (
	"fmt"
	"regexp"
)

type ReList []*regexp.Regexp

func (rl ReList) FindStringIndex(s string) []int {
	for _, re := range rl {
		if idx := re.FindStringIndex(s); idx != nil {
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
		`\*`, `/`, `%`, `\^`, `\?`, `:`,
	}
	reCMultiOps = []string{
		`(?:\S!|!\s)`,                   // !
		`(?:\s\+{2}\s)`, `(?:\s-{2}\s)`, // ++, --
		`(?:[^\s+]\+{1}[^+]|[^+]\+{1}[^\s+=])`,           // +
		`(?:[^\s-(]-[^-]|-[^\s\da-zA-Z=-])`,              // -
		`(?:[^\s=!><+-\\*/&|^%]=|=[^\s=])`,               // =
		`(?:[^\s><][=!><+-\\*/&|^%]=)`,                   // *=
		`(?:[^\s|]\||\|[^\s|=])`, `(?:[^\s&]&|&[^\s&=])`, // |, &
		`(?:\S[^& ]&|&[^& ]\S)`, `(?:\S[^| ]\||\|[^| ]\S)`, // &&, ||
		`(?:[^\s<]<|<[^\s<=])`, `(?:[^\s>]>|>[^\s>=])`, // <, <<, >, >>
	}
	ReCOperators ReList
)

func init() {
	for _, it := range CKeywords {
		re := fmt.Sprintf(`\b(?:%s)\b\S`, it)
		ReCKeywords = append(ReCKeywords, regexp.MustCompile(re))
	}

	for _, it := range reCSimpleOps {
		re := fmt.Sprintf(`(?:\S%s|%s[^\s=])`, it, it)
		ReCOperators = append(ReCOperators, regexp.MustCompile(re))
	}
	for _, it := range reCMultiOps {
		ReCOperators = append(ReCOperators, regexp.MustCompile(it))
	}
}
