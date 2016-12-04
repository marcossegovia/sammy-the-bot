package main

import (
	"log"
	"fmt"

	"github.com/spf13/viper"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
)

func main() {

	configurationFileHandler := viper.New()
	configurationFileHandler.SetConfigName("sammy_config")
	configurationFileHandler.AddConfigPath("$GOPATH/src/github.com/MarcosSegovia/SammyTheBot/config")
	err := configurationFileHandler.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	botToken := configurationFileHandler.GetString("configuration.sammyBotToken")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	sammyBrainFileHandler := viper.New()
	sammyBrainFileHandler.SetConfigName("sammy_brain")
	sammyBrainFileHandler.AddConfigPath("$GOPATH/src/github.com/MarcosSegovia/SammyTheBot/sammy")
	err = sammyBrainFileHandler.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.MessageConfig{}
		if "Hi" == update.Message.Text {
			some_salutations := sammyBrainFileHandler.GetStringSlice("welcome.salutations")
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, randomStringFromStringSlice(some_salutations))
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		}
		bot.Send(msg)
	}
}

func randomStringFromStringSlice(stringOfSlices []string) (string) {
	return stringOfSlices[rand.Intn(len(stringOfSlices))]
}

