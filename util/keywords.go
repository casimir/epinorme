package util

var (
	CKeywords = []string{
		"auto",
		"char",
		"const",
		"continue",
		"do",
		"double",
		"else",
		"enum",
		"extern",
		"float",
		"if",
		"int",
		"long",
		"register",
		"return",
		"short",
		"signed",
		"sizeof",
		"static",
		"struct",
		"typedef",
		"union",
		"unsigned",
		"void",
		"volatile",
		"while",
	}
	CKeywordsForbidden = []string{
		"break",
		"case",
		"default",
		"for",
		"goto",
		"switch",
	}

	COperators = []string{
		"+", "-", "*", "/", "%",
		"==", "!=", ">", "<", ">=", "<=",
		"&", "|", "^", "<<", ">>",
		"&&", "||", "!",
		"=", "+=", "-+", "*=", "/=", "%=", "<<=", "<<=", "&=", "^=", "!=",
		"?", ":",
	}
)
