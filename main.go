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

	bot := sammy.NewSammySpeaker(brain)
	token := cfg.GetString("configuration.token")
	//commands := setCommands(brain)
	api, err := tgbotapi.NewBotAPI(token)
	check(err, "could not initialize bot: %v")
	log.Printf("Authorized on account %v", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Printf("[%v] %v", update.Message.From.UserName, update.Message.Text)

		req := sammy.Request(update.Message.Text)
		resp := bot.Process(req)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.String())
		api.Send(msg)
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
func setCommands(brain *viper.Viper) map[string]string {
	var commands = map[string]string{}
	availableCmds := brain.GetStringSlice("commands.commands")
	for _, cmd := range availableCmds {
		index := cmd[1:]
		commands[index] = cmd
	}
	return commands
}
