package file

type Token struct {
	Type     TokenType
	Position int
	Content  string
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	WS
	IDENT
	COLON
	CONTINUATION
	FLAG
)

func (tt TokenType) String() string {
	switch tt {
	case ILLEGAL:
		return "illegal"
	case EOF:
		return "eof"
	case WS:
		return "whitespace"
	case IDENT:
		return "ident"
	case COLON:
		return "colon"
	case CONTINUATION:
		return "continuation"
	case FLAG:
		return "flag"
	}
	return ""
}
