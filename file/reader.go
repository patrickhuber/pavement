package file

import (
	"fmt"

	"github.com/patrickhuber/pavement/directive"
)

type reader struct {
	current *directive.Directive
	scanner Scanner
}

type Reader interface {
	Next() (bool, error)
	Current() *directive.Directive
}

func NewReader(s Scanner) Reader {
	return &reader{
		scanner: s,
	}
}

// Current implements Reader
func (r *reader) Current() *directive.Directive {
	return r.current
}

// Next implements Reader
func (r *reader) Next() (bool, error) {
	// read until we get a directive
	for {
		token, err := r.scanner.Scan()
		if err != nil {
			return false, r.err("%w", err)
		}
		switch token.Type {
		case EOF:
			return false, nil
		case ILLEGAL:
			return false, r.err("illegal token")
		case WS:
			continue
		case IDENT:
			directive, err := r.directive(token)
			if err != nil {
				return false, err
			}
			r.current = directive
			return true, nil
		default:
			return false, r.err("unrecognized token")
		}
	}
}

func (r *reader) directive(token *Token) (*directive.Directive, error) {
	dt, ok := directive.ParseType(token.Content)
	if !ok {
		return nil, r.err("invalid directive type '%s'", token.Content)
	}

	switch dt {
	case directive.FromType:
		return r.from()
	case directive.RunType:
		return r.run()
	default:
		return nil, r.err("unable to parse directive %s", token.Content)
	}
}

// FROM<WS>ubuntu:latest[WS]<EOL|EOF>
// FROM<WS>ubuntu[WS]<EOL|EOF>
func (r *reader) from() (*directive.Directive, error) {
	err := r.whitespace()
	if err != nil {
		return nil, err
	}

	t, err := r.expect(IDENT)
	if err != nil {
		return nil, err
	}
	base := t.Content

	version := ""
	_, ok, err := r.optional(COLON)
	if err != nil {
		return nil, err
	}
	if ok {
		t, err := r.expect(IDENT)
		if err != nil {
			return nil, err
		}
		version = t.Content
	}

	// optional whitespace
	_, ok, err = r.optional(WS)
	if err != nil {
		return nil, err
	}

	// read eol or eof if previous was ws
	if ok {
		_, err = r.any(EOF, EOL)
		if err != nil {
			return nil, err
		}
	}

	return &directive.Directive{
		Type: directive.FromType,
		From: &directive.From{
			Base:    base,
			Version: version,
		},
	}, nil
}

func (r *reader) run() (*directive.Directive, error) {
	command := ""

	_, err := r.expect(WS)
	if err != nil {
		return nil, err
	}

	t, err := r.expect(IDENT)
	if err != nil {
		return nil, err
	}
	command = t.Content

	args, err := r.args()
	if err != nil {
		return nil, err
	}

	return &directive.Directive{
		Type: directive.RunType,
		Run: &directive.Run{
			Command:   command,
			Arguments: args,
		},
	}, nil
}

func (r *reader) args() ([]string, error) {
	var arguments []string
	for {
		t, err := r.read()
		if err != nil {
			return nil, err
		}
		if t.Type == EOF {
			break
		}
		if t.Type == EOL {
			break
		}
		if t.Type == WS {
			continue // skip whitespace
		}
		if t.Type != IDENT && t.Type != FLAG && t.Type != STRING {
			return nil, r.err("expected IDENT, FLAG or STRING but found %s(%s)", t.Type.String(), t.Content)
		}
		arguments = append(arguments, t.Content)
	}
	return arguments, nil
}

func (r *reader) whitespace() error {
	t, err := r.read()
	if err != nil {
		return err
	}
	if t.Type != WS {
		return r.err("expected whitespace")
	}
	return nil
}

func (r *reader) expect(tokenType TokenType) (*Token, error) {
	t, err := r.read()
	if err != nil {
		return nil, err
	}
	if t.Type != tokenType {
		return nil, r.err("expected %s", t.Type.String())
	}
	return t, nil
}

func (r *reader) optional(tokenType TokenType) (*Token, bool, error) {
	t, err := r.read()
	if err != nil {
		return nil, false, err
	}
	ok := t.Type == tokenType
	return t, ok, nil
}

func (r *reader) any(tokenTypes ...TokenType) (*Token, error) {
	if len(tokenTypes) == 0 {
		return nil, r.err("no token types specified")
	}
	t, err := r.read()
	if err != nil {
		return nil, err
	}
	for _, tt := range tokenTypes {
		if t.Type == tt {
			return t, nil
		}
	}
	return nil, fmt.Errorf("token type '%s' does not match any required type", t.Type.String())
}

// read reads the token and if the token is a CONTINUATION it replaces the CONTINUATION and EOL with an empty WS token
func (r *reader) read() (*Token, error) {

	t, err := r.scanner.Scan()
	if err != nil {
		return nil, err
	}

	// if the token is not a continuation, we are done
	if t.Type != CONTINUATION {
		return t, nil
	}

	// read the next token.
	t, err = r.scanner.Scan()
	if err != nil {
		return nil, err
	}

	//  we are expecting a EOL token
	if t.Type != EOL {
		return nil, r.err("EOL expected after CONTINUATION found '%s'", t.Type.String())
	}

	// return empty whitespace
	return &Token{Type: WS, Position: r.scanner.Position()}, nil
}

func (r *reader) err(message string, param ...any) error {
	return NewReaderError(r.scanner, message, param...)
}

type readerError struct {
	pos     int
	message string
	param   []any
}

func NewReaderError(scanner Scanner, message string, param ...any) error {
	return &readerError{pos: scanner.Position(), message: message, param: param}
}

func (err *readerError) Error() string {
	var parameters []any
	parameters = append(parameters, err.pos)
	parameters = append(parameters, err.message)
	parameters = append(parameters, err.param...)
	return fmt.Sprintf("error at position %d: %s", parameters...)
}
