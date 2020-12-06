package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/marcossegovia/sammy-the-bot/help"
	"github.com/marcossegovia/sammy-the-bot/sammy"
	"github.com/marcossegovia/sammy-the-bot/start"
	"github.com/marcossegovia/sammy-the-bot/weather"
	"github.com/spf13/viper"
)

func main() {
	brain, err := read("sammy_brain")
	check(err, "could not read config file: %v")

	token := brain.GetString("configuration.telegram")
	api, err := tgbotapi.NewBotAPI(token)
	check(err, "could not initialize bot: %v")
	log.Printf("Authorized on account %v", api.Self.UserName)

	s := sammy.NewSammy(brain, api)

	var commands *[]sammy.Command
	commands = loadCommands(s)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.Api.GetUpdatesChan(u)
	check(err, "could not get telegram message: %v")
	for update := range updates {
		msg := update.Message
		if update.CallbackQuery != nil {
			s.Api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			s.Api.Send(msg)
			check(err, "could not send message because: %v")
			continue
		}
		if update.Message == nil {
			continue
		}
		for _, cmd := range *commands {
			_, err := cmd.Evaluate(msg)
			check(err, "command failed: %v")
		}
	}
}
func loadCommands(s *sammy.Sammy) *[]sammy.Command {
	cmds := new([]sammy.Command)
	cnames := []string{}

	startCmd := start.NewStart(s)
	weatherCmd := weather.NewWeather(s)

	*cmds = append(*cmds, startCmd)
	*cmds = append(*cmds, weatherCmd)

	cnames = append(cnames, weatherCmd.Description())
	*cmds = append(*cmds, help.NewHelp(s, cnames))

	return cmds
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
