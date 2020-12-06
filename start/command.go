package start

import (
	"bytes"
	"fmt"
	"github.com/marcossegovia/sammy-the-bot/sammy"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Start struct {
	sammy *sammy.Sammy
	cmd   *sammy.Cmd
}

func NewStart(sam *sammy.Sammy) *Start {
	s := new(Start)
	s.sammy = sam
	s.cmd = sammy.NewCommand("start", "/start", "Initialize Sammy :D")
	return s
}

func (s *Start) Evaluate(msg *tgbotapi.Message) (bool, error) {
	if msg.Text != s.cmd.Exec {
		return false, nil
	}

	var buffer bytes.Buffer
	buffer.WriteString("Hi there !\n")
	buffer.WriteString("Im your botpher assistance on whatever you need.\n My source code is in https://github.com/marcossegovia/sammy-the-bot\n Just follow /help to see things I can do. \n\n")
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	_, err := s.sammy.Api.Send(newMsg)
	if err != nil {
		return false, fmt.Errorf("could not send message because: %v", err)
	}
	return true, nil
}

func (s *Start) Description() string {
	return s.cmd.Exec + " - " + s.cmd.Desc
}
