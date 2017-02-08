package command

import "bytes"

type Command interface {
	Description() string
	Evaluate() bytes.Buffer
}

type Cmd struct {
	Tag  string
	Exec string
	Desc string
}

func NewCommand(tag, exec, desc string) *Cmd {
	c := new(Cmd)
	c.Tag = tag
	c.Exec = exec
	c.Desc = desc
	return c
}
