package main

import (
	"log"
	"math/rand"

	"github.com/spf13/viper"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"bytes"
)

var cfg, brain *viper.Viper

func init() {
	cfg = viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigName("sammy_config")
	err := cfg.ReadInConfig()
	if err != nil {
		log.Printf("could not read config file: %v", err)
	}
	brain = viper.New()
	brain.AddConfigPath(".")
	brain.SetConfigName("sammy_brain")
	err = brain.ReadInConfig()
	if err != nil {
		log.Printf("could not read config file: %v", err)
	}
}

func main() {
	token := cfg.GetString("configuration.token")
	var commands = map[string]string{}
	availableCmds := brain.GetStringSlice("commands.commands")
	for _, cmd := range availableCmds {
		index := cmd[1:]
		commands[index] = cmd
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("could not initialize bot: %v", err)
	}
	log.Printf("Authorized on account %v", bot.Self.UserName)

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%v] %v", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I do not know what to tell you.")
		if update.Message.IsCommand() {
			if name, ok := commands[update.Message.Text[1:]]; ok {
				var buffer bytes.Buffer
				buffer.WriteString("You asked me to do ")
				buffer.WriteString(name)
				buffer.WriteString(" and I can do it :)")
				msg.Text = buffer.String()
			}
		}
		if "Hi" == update.Message.Text {
			salutations := brain.GetStringSlice("welcome.salutations")
			msg.Text = salute(salutations)
		}
		log.Printf("I'm responding: %v", msg.Text)
		bot.Send(msg)
	}
}

func salute(salutations []string) (string) {
	return salutations[rand.Intn(len(salutations))]
}
