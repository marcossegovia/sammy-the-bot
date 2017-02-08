package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/MarcosSegovia/sammy-the-bot/sammy"
)

func main() {
	cfg, err := read("sammy_config")
	check(err, "could not read config file: %v")
	brain, err := read("sammy_brain")
	check(err, "could not read config file: %v")

	bot := sammy.NewSammySpeaker(brain, cfg)
	token := cfg.GetString("configuration.token")
	api, err := tgbotapi.NewBotAPI(token)
	check(err, "could not initialize bot: %v")

	log.Printf("Authorized on account %v", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := api.GetUpdatesChan(u)
	for update := range updates {
		//if update.CallbackQuery != nil {
		//	api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
		//	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		//	api.Send(msg)
		//	check(err, "could not send message because: %v")
		//	continue
		//}
		if update.Message == nil {
			continue
		}
		log.Printf("[%v] %v", update.Message.From.UserName, update.Message.Text)

		req := sammy.Request(update.Message.Text)
		resp := bot.Process(req)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.String())
		_, err := api.Send(msg)
		check(err, "could not send message because: %v")
	}
}

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}

func read(path string) (*viper.Viper, error) {
	f := viper.New()
	f.AddConfigPath(".")
	f.SetConfigName(path)
	err := f.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return f, nil
}
