package pavementfile

type File struct {
	From From
}

type From struct {
	Image string
	Tag   string
}
