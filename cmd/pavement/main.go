package main

import (
	"github.com/alecthomas/kong"
	"github.com/patrickhuber/pavement/internal/docker"
)

var cli struct {
	Apply ApplyCommand `cmd:"" help:"Remove files."`
}

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "apply <path>":
		err := ctx.Run(&Context{})
		ctx.FatalIfErrorf(err)
	default:
		panic(ctx.Command())
	}
}

type Context struct {
	Debug bool
}

type ApplyCommand struct {
	Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
}

func (r *ApplyCommand) Run(ctx *Context) error {
	provider := &docker.Provider{}
	return provider.Run()
}
