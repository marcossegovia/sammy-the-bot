package main

import (
	"log"
	"net/http"

	"github.com/marcossegovia/sammy-the-bot/help"
	"github.com/marcossegovia/sammy-the-bot/sammy"
	"github.com/marcossegovia/sammy-the-bot/start"
	"github.com/marcossegovia/sammy-the-bot/weather"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/marcossegovia/sammy-the-bot/github"
	"github.com/marcossegovia/sammy-the-bot/user"
)

func main() {
	brain, err := read("sammy_brain")
	check(err, "could not read config file: %v")

	token := brain.GetString("configuration.telegram")
	api, err := tgbotapi.NewBotAPI(token)
	check(err, "could not initialize bot: %v")
	log.Printf("Authorized on account %v", api.Self.UserName)

	userRepository := user.NewUserRepository("sammy-host:6379", "", 0)

	s := sammy.NewSammy(brain, api, userRepository)

	hook := github.NewHook(s)
	http.Handle("/github/hooks/", hook)
	go http.ListenAndServe(":80", nil)

	var commands *[]sammy.Command
	commands = loadCommands(s)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.Api.GetUpdatesChan(u)
	check(err, "could not get telegram message: %v")
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
		commanded := false
		for _, cmd := range *commands {
			eval, err := cmd.Evaluate(msg)
			if eval {
				commanded = true
			}
			check(err, "command failed: %v")
		}
		if !commanded {
			err = s.Process(msg)
			check(err, "conversation failed: %v")
		}
		commanded = false
	}
}
func loadCommands(s *sammy.Sammy) *[]sammy.Command {
	cmds := new([]sammy.Command)
	cnames := []string{}

	startCmd := start.NewStart(s)
	*cmds = append(*cmds, startCmd)
	weatherCmd := weather.NewWeather(s)
	*cmds = append(*cmds, weatherCmd)
	githubCmd := github.NewGithub(s)
	*cmds = append(*cmds, githubCmd)
	cnames = append(cnames, weatherCmd.Description())
	cnames = append(cnames, githubCmd.Description())
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
