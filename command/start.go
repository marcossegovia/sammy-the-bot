package command

import "bytes"

type Start struct {
	Cmd *Cmd
}

func NewStart() *Start {
	s := new(Start)
	s.Cmd = NewCommand("start", "/start", "Initialize Sammy :D")
	return s
}

func (s *Start) Evaluate() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("Im your botpher assistance on whatever you need.\n My source code is in https://github.com/MarcosSegovia/sammy-the-bot\n Just follow /help to see things I can do. \n\n")
	return buffer
}

func (s *Start) Description() string {
	return s.Cmd.Exec + " - " + s.Cmd.Desc
}
