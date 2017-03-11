package sammy

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/marcossegovia/sammy-the-bot/user"
	"github.com/satori/go.uuid"
)

type Response struct {
	Response string
	Status   int
}

func (r Response) String() string {
	return string(r.Response)
}

func NewUser(chatId int64, userName string) *user.User {
	return &user.User{Id: uuid.NewV4().String(), ChatId: chatId, Name: userName}
}

type Sammy struct {
	Brain  *viper.Viper
	Config *viper.Viper
	Api    *tgbotapi.BotAPI
	users  *user.UserRepository
}

func NewSammy(brain *viper.Viper, api *tgbotapi.BotAPI, userRepository *user.UserRepository) *Sammy {
	s := new(Sammy)
	s.Brain = brain
	s.Api = api
	s.users = userRepository
	return s
}

func (s *Sammy) AddUser(user *user.User) error {
	return s.users.AddUser(user)
}

func (s *Sammy) GetUser(userId string) (*user.User, error) {
	return s.users.GetUser(userId)
}

func (s *Sammy) GetUserIdByChatId(chatId int64) (string, error){
	return s.users.GetUserId(chatId)
}
