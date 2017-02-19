package sammy

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

type Response struct {
	Response string
	Status   int
}

func (r Response) String() string {
	return string(r.Response)
}

type User struct {
	ChatId int64
	Name   string
}

type Sammy struct {
	Brain    *viper.Viper
	Config   *viper.Viper
	Api      *tgbotapi.BotAPI
	userChat map[int64]*User
}

func NewSammy(brain *viper.Viper, api *tgbotapi.BotAPI) *Sammy {
	s := new(Sammy)
	s.Brain = brain
	s.Api = api
	s.userChat = make(map[int64]*User)
	return s
}

func (s *Sammy) AddChatId(chatId int64, username string) {
	s.userChat[chatId] = &User{chatId, username}
}

func (s *Sammy) GetUser(chatId int64) (*User, error) {
	if user, ok := s.userChat[chatId]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("%v", "there is not a user registered given the chatId")
}
