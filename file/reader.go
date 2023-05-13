package file

import (
	"fmt"
	"strings"
)

type reader struct {
	current *Directive
	scanner Scanner
}

type Reader interface {
	Next() (bool, error)
	Current() *Directive
}

func NewReader(s Scanner) Reader {
	return &reader{
		scanner: s,
	}
}

// Current implements Reader
func (r *reader) Current() *Directive {
	return r.current
}

// Next implements Reader
func (r *reader) Next() (bool, error) {
	// read until we get a directive
	for {
		token, err := r.scanner.Scan()
		if err != nil {
			return false, err
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

func (r *reader) directive(token *Token) (*Directive, error) {
	dt, ok := ParseDirectiveType(token.Content)
	if !ok {
		return nil, r.err("invalid directive type '%s'", token.Content)
	}
	switch dt {
	case From:
		return r.from()
	case Run:
		return r.run()
	default:
		return nil, r.err("unable to parse directive %s", token.Content)
	}
}

func (r *reader) from() (*Directive, error) {
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
	return &Directive{
		Type: From,
		From: &FromDirective{
			Base:    base,
			Version: version,
		},
	}, nil
}

func (r *reader) run() (*Directive, error) {
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

	return &Directive{
		Type: Run,
		Run: &RunDirective{
			Command:   command,
			Arguments: args,
		},
	}, nil
}

func (r *reader) args() ([]string, error) {
	var arguments []string
	for {
		t, err := r.scanner.Scan()
		if err != nil {
			return nil, err
		}
		if t.Type == EOF {
			break
		}
		if t.Type == CONTINUATION {
			continue // skip continuation
		}
		if t.Type == WS {
			if strings.ContainsRune(t.Content, '\r') {
				break // break on newline
			}
			continue // skip whitespace
		}
		if t.Type != IDENT && t.Type != FLAG {
			return nil, r.err("expected IDENT or FLAG, found '%s'", t.Content)
		}
		arguments = append(arguments, t.Content)
	}
	return arguments, nil
}

func (r *reader) whitespace() error {
	t, err := r.scanner.Scan()
	if err != nil {
		return err
	}
	if t.Type != WS {
		return r.err("expected whitespace")
	}
	return nil
}

func (r *reader) expect(tokenType TokenType) (*Token, error) {
	t, err := r.scanner.Scan()
	if err != nil {
		return nil, err
	}
	if t.Type != tokenType {
		return nil, r.err("expected %s", t.Type.String())
	}
	return t, nil
}

func (r *reader) optional(tokenType TokenType) (*Token, bool, error) {
	t, err := r.scanner.Scan()
	if err != nil {
		return nil, false, err
	}
	ok := t.Type == tokenType
	return t, ok, nil
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
