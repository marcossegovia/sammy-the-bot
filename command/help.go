package command

import "bytes"

type Help struct {
	Cmd    *Cmd
	cnames []string
}

func NewHelp(cnames []string) *Help {
	h := new(Help)
	h.Cmd = NewCommand("help", "/help", "Show available commands")
	h.cnames = cnames
	return h
}

func (h *Help) Evaluate() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("Here is what I can do: \n\n")
	for _, cmd := range h.cnames {
		buffer.WriteString(cmd)
	}
	return buffer
}

func (h *Help) Description() string {
	return h.Cmd.Exec + " - " + h.Cmd.Desc
}
