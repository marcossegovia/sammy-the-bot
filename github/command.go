package github

import (
	"bytes"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/marcossegovia/sammy-the-bot/sammy"
)

type Github struct {
	sammy *sammy.Sammy
	cmd   *sammy.Cmd
}

func NewGithub(sam *sammy.Sammy) *Github {
	s := new(Github)
	s.sammy = sam
	s.cmd = sammy.NewCommand("github", "/github", "Get your webhook url to integrate github")
	return s
}

func (g *Github) Evaluate(msg *tgbotapi.Message) (bool, error) {
	if msg.Text != g.cmd.Exec {
		return false, nil
	}
	userId, err := g.sammy.GetUserIdByChatId(msg.Chat.ID)
	if err != nil {
		return false, fmt.Errorf("could not get user id because: %v", err)
	}
	if userId == "" && err == nil {
		return false, fmt.Errorf("could not get user id because: %v", fmt.Errorf("%v", "there is not a user registered given the chatId"))
	}

	var buffer bytes.Buffer
	buffer.WriteString("In order to integrate sammy with github, you should add this url to your repository webhooks http://138.68.68.19/github/hooks/" + userId + "\n")
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	_, err = g.sammy.Api.Send(newMsg)
	if err != nil {
		return false, fmt.Errorf("could not send message because: %v", err)
	}
	return true, nil
}

func (g *Github) Description() string {
	return g.cmd.Exec + " - " + g.cmd.Desc
}
