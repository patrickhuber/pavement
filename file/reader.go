package file

type reader struct {
	current *Directive
	scanner Scanner
}

type Reader interface {
	Next() bool
	Current() *Directive
}

func NewReader(s Scanner) Reader {
	return &reader{
		scanner: s,
	}
}

// Current implements Parser
func (r *reader) Current() *Directive {
	return r.current
}

// Next implements Parser
func (r *reader) Next() bool {
	return false
}
