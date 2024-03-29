package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

type scanner struct {
	reader   *bufio.Reader
	position int
	column   int
}

type Scanner interface {
	Scan() (*Token, error)
	Position() int
}

func NewScanner(reader io.Reader) Scanner {
	return &scanner{
		reader:   bufio.NewReader(reader),
		position: 0,
		column:   0,
	}
}

func (s *scanner) Position() int {
	return s.position
}

func (s *scanner) Scan() (*Token, error) {
	r, err := s.read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return s.token(EOF, ""), nil
		}
		return nil, err
	}
	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	switch {
	case r == '\n' || r == '\r':
		if err := s.unread(); err != nil {
			return nil, err
		}
		return s.endOfLine()
	case unicode.IsSpace(r):
		if err := s.unread(); err != nil {
			return nil, err
		}
		return s.whitespace()
	case unicode.IsLetter(r):
		if err := s.unread(); err != nil {
			return nil, err
		}
		return s.identifier()
	case r == '\\':
		return s.token(CONTINUATION, string(r)), nil
	case r == ':':
		return s.token(COLON, string(r)), nil
	case r == '-' || r == '/':
		if err := s.unread(); err != nil {
			return nil, err
		}
		return s.flag()
	case r == '"' || r == '\'':
		if err := s.unread(); err != nil {
			return nil, err
		}
		return s.string()
	}

	return &Token{Type: ILLEGAL, Position: s.position, Content: string(r)}, nil
}

func (s *scanner) string() (*Token, error) {
	var start rune
	var buf bytes.Buffer
	for i := 0; true; i++ {

		r, err := s.read()
		if err != nil {
			return nil, err
		}

		buf.WriteRune(r)

		if i == 0 {
			start = r
			continue
		}

		if start == r {
			break
		}
	}
	return s.token(STRING, buf.String()), nil
}

func (s *scanner) endOfLine() (*Token, error) {

	var buf bytes.Buffer
	for i := 0; i < 2; i++ {

		// read the next character
		r, err := s.read()
		if err != nil {
			return nil, err
		}

		// if we don't recognize the character, bail
		switch r {
		case '\n':
			buf.WriteRune(r)
			return s.token(EOL, buf.String()), nil
		case '\r':
			buf.WriteRune(r)
		default:
			return nil, fmt.Errorf("unrecognized newline character %c", r)
		}
	}
	return nil, fmt.Errorf("unable to read end of line characters")
}

func (s *scanner) whitespace() (*Token, error) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	for {
		r, err := s.read()
		switch {
		case err != nil:
			if errors.Is(io.EOF, err) {
				return s.token(WS, buf.String()), nil
			}
			return nil, err
		case r == '\t' || r == '\v' || r == '\f' || r == ' ':
			if _, err = buf.WriteRune(r); err != nil {
				return nil, err
			}
		default:
			if err = s.unread(); err != nil {
				return nil, err
			}
			return s.token(WS, buf.String()), nil
		}
	}
}

func (s *scanner) identifier() (*Token, error) {
	var buf bytes.Buffer
	i := 0
	for {
		r, err := s.read()
		switch {
		case err != nil:
			if errors.Is(io.EOF, err) {
				return s.token(IDENT, buf.String()), nil
			}
			return nil, err
		case unicode.IsLetter(r) || (unicode.IsNumber(r) && i > 0):
			if _, err = buf.WriteRune(r); err != nil {
				return nil, err
			}
		default:
			if err = s.unread(); err != nil {
				return nil, err
			}
			return s.token(IDENT, buf.String()), nil
		}
		i++
	}
}

func (s *scanner) flag() (*Token, error) {
	var buf bytes.Buffer
	i := 0
	for {
		r, err := s.read()
		switch {
		case err != nil:
			if errors.Is(io.EOF, err) {
				return s.token(FLAG, buf.String()), nil
			}
			return nil, err
		case unicode.IsLetter(r) || (unicode.IsNumber(r) && i > 0) || r == '/' && i == 0 || r == '-' && i <= 1:
			if _, err = buf.WriteRune(r); err != nil {
				return nil, err
			}
		default:
			if err = s.unread(); err != nil {
				return nil, err
			}
			return s.token(FLAG, buf.String()), nil
		}
		i++
	}
}

func (s *scanner) token(t TokenType, c string) *Token {
	return &Token{
		Type:     t,
		Position: s.position - len(c),
		Content:  c,
	}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *scanner) read() (rune, error) {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return rune(0), err
	}
	s.position++
	return ch, nil
}

// unread places the previously read rune back on the reader.
func (s *scanner) unread() error {
	s.position--
	return s.reader.UnreadRune()
}
