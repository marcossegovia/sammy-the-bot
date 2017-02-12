package sammy

import (
	"bytes"
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	NO_RESPONSE  = 0
	CONVERSATION = 1
)

func (s *Sammy) Process(msg *tgbotapi.Message) error {
	log.Printf("[%v] %v", msg.From.UserName, msg.Text)
	buffer := bytes.Buffer{}
	resp := Response{"I do not know what to tell you. Maybe you need /help", NO_RESPONSE}
	if "Hi" == msg.Text {
		salutations := s.Brain.GetStringSlice("welcome.salutations")
		resp = Response{salute(salutations), CONVERSATION}
	}
	log.Printf("I'm responding: %v", resp)

	buffer.WriteString(resp.String())
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	_, err := s.Api.Send(newMsg)
	if err != nil {
		return fmt.Errorf("could not send message because: %v", err)
	}
	return nil
}
