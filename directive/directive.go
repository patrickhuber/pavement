package directive

type Directive struct {
	Type DirectiveType
	From *From
	Run  *Run
}

type DirectiveType int

const (
	FromType DirectiveType = iota
	RunType
)

func (dt DirectiveType) String() string {
	switch dt {
	case FromType:
		return "FROM"
	case RunType:
		return "RUN"
	default:
		return "UNKNOWN"
	}
}

func ParseType(s string) (DirectiveType, bool) {
	switch s {
	case "FROM":
		return FromType, true
	case "RUN":
		return RunType, true
	}
	return DirectiveType(-1), false
}

type From struct {
	Base    string
	Version string
}

type Run struct {
	Command   string
	Arguments []string
}
