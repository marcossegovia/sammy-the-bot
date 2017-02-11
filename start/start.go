package start

import (
	"bytes"
	"log"

	"github.com/MarcosSegovia/sammy-the-bot/sammy"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Start struct {
	sammy *sammy.Sammy
	Cmd   *sammy.Cmd
}

func NewStart(sam *sammy.Sammy) *Start {
	s := new(Start)
	s.sammy = sam
	s.Cmd = sammy.NewCommand("start", "/start", "Initialize Sammy :D")
	return s
}

func (s *Start) Evaluate(msg *tgbotapi.Message) {
	if msg.Text != s.Cmd.Exec {
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("Im your botpher assistance on whatever you need.\n My source code is in https://github.com/MarcosSegovia/sammy-the-bot\n Just follow /help to see things I can do. \n\n")
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	_, err := s.sammy.Api.Send(newMsg)
	check(err, "could not send message because: %v")
}

func (s *Start) Description() string {
	return s.Cmd.Exec + " - " + s.Cmd.Desc
}

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
