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
	buffer.WriteString("You asked me to do ")
	buffer.WriteString(c.Exec)
	buffer.WriteString(" and I can do it :)")
	return buffer
}

func NewCommand(tag, exec string) *Cmd {
	c := new(Cmd)
	c.Tag = tag
	c.Exec = exec
	return c
}
