package sammy

import (
	"bytes"
	"fmt"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/marcossegovia/apiai-go"
)

type Ai struct {
	sammy       *Sammy
	apiaiClient *apiai.ApiClient
}

func NewAiConversation(s *Sammy) (*Ai, error) {
	token := s.Brain.GetString("configuration.api_ai_client")
	client, err := apiai.NewClient(
		&apiai.ClientConfig{
			Token:      token,
			SessionId:  "faf7b9eb-7702-4b6b-bd43-d58147d156c4",
			QueryLang:  "en",    //Default en
			SpeechLang: "en-US", //Default en-US
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &Ai{sammy: s, apiaiClient: client}, nil
}

func (ai *Ai) Query(message string) (string, error) {
	qr, err := ai.apiaiClient.Query(apiai.Query{Query: []string{message}})
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	return qr.Result.Fulfillment.Speech, nil
}

func (s *Sammy) Process(msg *tgbotapi.Message) error {
	log.Printf("[%v] %v", msg.From.UserName, msg.Text)
	ai, err := NewAiConversation(s)
	if err != nil {
		return fmt.Errorf("could not initialize apiai, because: %v", err)
	}
	response, err := ai.Query(msg.Text)
	if err != nil {
		return fmt.Errorf("error on getting an response, because: %v", err)
	}
	log.Printf("I'm responding: %v", response)

	buffer := bytes.Buffer{}
	buffer.WriteString(response)
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, buffer.String())
	_, err = s.Api.Send(newMsg)
	if err != nil {
		return fmt.Errorf("could not send message because: %v", err)
	}
	return nil
}
