package command

import "bytes"

type Command interface {
	Evaluate() bytes.Buffer
}

type Cmd struct {
	Tag  string
	Exec string
}

func (c *Cmd) Evaluate() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("I have no action on ")
	buffer.WriteString(c.Exec)
	return buffer
}

func NewCommand(tag, exec string) *Cmd {
	c := new(Cmd)
	c.Tag = tag
	c.Exec = exec
	return c
}
