package main

import (
	"log"

	"github.com/spf13/viper"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/MarcosSegovia/sammy-the-bot/sammy"
	"github.com/MarcosSegovia/sammy-the-bot/start"
	"github.com/MarcosSegovia/sammy-the-bot/weather"
	"github.com/MarcosSegovia/sammy-the-bot/help"
)

func main() {
	var commands *[]sammy.Command
	brain, err := read("sammy_brain")
	check(err, "could not read config file: %v")

	token := brain.GetString("configuration.telegram")
	api, err := tgbotapi.NewBotAPI(token)
	check(err, "could not initialize bot: %v")
	log.Printf("Authorized on account %v", api.Self.UserName)

	s := sammy.NewSammy(brain, api)
	commands = loadCommands(s)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.Api.GetUpdatesChan(u)
	for update := range updates {
		msg := update.Message
		if update.CallbackQuery != nil {
			//	sammy.Api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
			//	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			//	sammy.Api.Send(msg)
			//	check(err, "could not send message because: %v")
			//	continue
		}
		if update.Message == nil {
			continue
		}
		log.Printf("[%v] %v", update.Message.From.UserName, update.Message.Text)

		//resp := Response{"I do not know what to tell you. Maybe you need /help", NO_RESPONSE}
		//if "Hi" == msg.Text {
		//	salutations := sammy.Brain.GetStringSlice("welcome.salutations")
		//	resp = Response{salute(salutations), CONVERSATION}
		//}
		for _, cmd := range *commands {
			cmd.Evaluate(msg)
		}
		//log.Printf("I'm responding: %v", resp)
	}
}
func loadCommands(s *sammy.Sammy) (*[]sammy.Command) {
	cmds := new([]sammy.Command)
	cnames := []string{}

	startCmd := start.NewStart(s)
	*cmds = append(*cmds, startCmd)
	weatherCmd := weather.NewWeather(s)
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
