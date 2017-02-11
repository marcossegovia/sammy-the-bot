package sammy

import "github.com/go-telegram-bot-api/telegram-bot-api"

type Command interface {
	Description() string
	Evaluate(msg *tgbotapi.Message)
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
