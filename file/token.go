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
)
