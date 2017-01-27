package main

import (
	"log"
	"math/rand"

	"github.com/spf13/viper"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"bytes"
)

func main() {
	cfg, err := read("sammy_config")
	if err != nil {
		log.Printf("could not read config file: %v", err)
	}
	brain, err := read("sammy_brain")
	if err != nil {
		log.Printf("could not read config file: %v", err)
	}
	token := cfg.GetString("configuration.token")
	commands := setCommands(brain)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("could not initialize bot: %v", err)
	}
	log.Printf("Authorized on account %v", bot.Self.UserName)

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
func setCommands(brain *viper.Viper) map[string]string {
	var commands = map[string]string{}
	availableCmds := brain.GetStringSlice("commands.commands")
	for _, cmd := range availableCmds {
		index := cmd[1:]
		commands[index] = cmd
	}
	return commands
}

func salute(salutations []string) (string) {
	return salutations[rand.Intn(len(salutations))]
}
