package file

type Directive struct {
	Type DirectiveType
	From *FromDirective
	Run  *RunDirective
}

type DirectiveType int

const (
	From DirectiveType = iota
	Run
)

func (dt DirectiveType) String() string {
	switch dt {
	case From:
		return "FROM"
	case Run:
		return "RUN"
	default:
		return "UNKNOWN"
	}
}

func ParseDirectiveType(s string) (DirectiveType, bool) {
	switch s {
	case "FROM":
		return From, true
	case "RUN":
		return Run, true
	}
	return DirectiveType(-1), false
}

type FromDirective struct {
	Base    string
	Version string
}

type RunDirective struct {
	Command   string
	Arguments []string
}
