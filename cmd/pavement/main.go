package main

import "github.com/alecthomas/kong"

var CLI struct {
	Apply ApplyCommand `cmd:"" help:"Remove files."`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "apply <path>":
	default:
		panic(ctx.Command())
	}
}

type Context struct{}

type ApplyCommand struct {
	Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
}

func (r *ApplyCommand) Run(ctx *Context) error {
	return nil
}
