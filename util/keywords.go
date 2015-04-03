package util

var (
	CKeywords = []string{
		"auto",
		"break",
		"case",
		"char",
		"const",
		"continue",
		"default",
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
		"switch",
		"typedef",
		"union",
		"unsigned",
		"void",
		"volatile",
		"while",
	}
	CKeywordsForbidden = []string{
		"for",
		"goto",
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
