package help

import (
	"bytes"
	"fmt"

	"github.com/MarcosSegovia/sammy-the-bot/sammy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Help struct {
	sammy  *sammy.Sammy
	cmd    *sammy.Cmd
	cnames []string
}

func NewHelp(s *sammy.Sammy, cnames []string) *Help {
	h := new(Help)
	h.sammy = s
	h.cmd = sammy.NewCommand("help", "/help", "Show available commands")
	h.cnames = cnames
	return h
}

func (h *Help) Evaluate(msg *tgbotapi.Message) (bool, error) {
	if msg.Text != h.cmd.Exec {
		return false, nil
	}
	var buffer bytes.Buffer
	buffer.WriteString("Here is what I can do: \n\n")
	buffer.WriteString("*Conversation*\n")
	buffer.WriteString("Say _Hi_ !\n\n")
	buffer.WriteString("*Commands*\n")
	for _, cmd := range h.cnames {
		buffer.WriteString(cmd)
	}
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	newMsg.ParseMode = "Markdown"
	_, err := h.sammy.Api.Send(newMsg)
	if err != nil {
		return false, fmt.Errorf("could not send message because: %v", err)
	}
	return true, nil
}

func (h *Help) Description() string {
	return h.cmd.Exec + " - " + h.cmd.Desc
}
