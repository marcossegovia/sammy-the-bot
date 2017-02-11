package sammy

import (
	"math/rand"

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

type Sammy struct {
	Brain  *viper.Viper
	Config *viper.Viper
	Api    *tgbotapi.BotAPI
}

func NewSammy(brain *viper.Viper, api *tgbotapi.BotAPI) *Sammy {
	s := new(Sammy)
	s.Brain = brain
	s.Api = api
	return s
}

func salute(salutations []string) string {
	return salutations[rand.Intn(len(salutations))]
}
