package help

import (
	"bytes"
	"log"

	"github.com/MarcosSegovia/sammy-the-bot/sammy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Help struct {
	sammy  *sammy.Sammy
	Cmd    *sammy.Cmd
	cnames []string
}

func NewHelp(s *sammy.Sammy, cnames []string) *Help {
	h := new(Help)
	h.sammy = s
	h.Cmd = sammy.NewCommand("help", "/help", "Show available commands")
	h.cnames = cnames
	return h
}

func (h *Help) Evaluate(msg *tgbotapi.Message) {
	if msg.Text != h.Cmd.Exec {
		return
	}
	var buffer bytes.Buffer
	buffer.WriteString("Here is what I can do: \n\n")
	for _, cmd := range h.cnames {
		buffer.WriteString(cmd)
	}
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	_, err := h.sammy.Api.Send(newMsg)
	check(err, "could not send message because: %v")
}

func (h *Help) Description() string {
	return h.Cmd.Exec + " - " + h.Cmd.Desc
}

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
