package file

type Directive struct {
	From *From
}

type From struct {
	Base    string
	Version string
}

type Run struct {
	Command   string
	Arguments []string
}
