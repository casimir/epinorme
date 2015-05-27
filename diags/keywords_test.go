package diags

import (
	"fmt"
	"testing"
)

var (
	operators2 = []string{
		"=", "+", "*", "/", "%", "&", "|", "^", "?", ":",
		"==", "!=", "<=", ">=",
		"+=", "-=", "*=", "/=", "%=", "&=", "|=",
		"&&", "||", "<<", ">>",
		"<<=", ">>=",
	}
	operatorsIncr = []string{
		"++", "--",
	}

	operatorsData = []struct {
		format   string
		expected bool
	}{
		{"a %s b", false},
		{"a%s b", true},
		{"a %sb", true},
	}
)

type TOp struct {
	str      string
	expected bool
}

func TestOperators(t *testing.T) {
	var tt []TOp
	for _, t := range operatorsData {
		for _, o := range operators2 {
			tt = append(tt, TOp{fmt.Sprintf(t.format, o), t.expected})
		}
		for _, o := range operatorsIncr {
			tt = append(tt, TOp{fmt.Sprintf(t.format, o), !t.expected})
		}
	}

	for _, it := range tt {
		got := ReCOperators.MatchString(it.str)
		if got != it.expected {
			t.Errorf("%q -> %t, want %t", it.str, got, it.expected)
		}
	}
}
