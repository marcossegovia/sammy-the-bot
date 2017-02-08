package sammy

import (
	"math/rand"
	"log"

	"github.com/spf13/viper"
	"github.com/MarcosSegovia/sammy-the-bot/command"
)

const (
	NO_RESPONSE  = 0
	CONVERSATION = 1
	COMMAND      = 2
)

type Request string
type Response struct {
	Response string
	Status   int
}

func (r Response) String() string {
	return string(r.Response)
}

type Sammy struct {
	brain    *viper.Viper
	config   *viper.Viper
	commands *[]interface{}
}

func NewSammySpeaker(brain, cfg *viper.Viper) *Sammy {
	s := new(Sammy)
	s.brain = brain
	s.config = cfg
	s.load()
	return s
}

func (sammy *Sammy) Process(req Request) Response {
	resp := Response{"I do not know what to tell you. Maybe you need /help", NO_RESPONSE}
	if "Hi" == req {
		salutations := sammy.brain.GetStringSlice("welcome.salutations")
		resp = Response{salute(salutations), CONVERSATION}
	}
	for _, v := range *sammy.commands {
		switch cmd := v.(type) {
		case *command.Start:
			if string(cmd.Cmd.Exec) == string(req) {
				resp = sammy.ProcessCmd(cmd)
			}
		case *command.Help:
			if string(cmd.Cmd.Exec) == string(req) {
				resp = sammy.ProcessCmd(cmd)
			}
		case *command.Weather:
			if string(cmd.Cmd.Exec) == string(req) {
				resp = sammy.ProcessCmd(cmd)
			}
		}
	}
	log.Printf("I'm responding: %v", resp)
	return resp
}

func (sammy *Sammy) ProcessCmd(cmd command.Command) Response {
	buffer := cmd.Evaluate()
	return Response{buffer.String(), COMMAND}
}

func (sammy *Sammy) load() {
	cmds := new([]interface{})
	cnames := []string{}

	startCmd := command.NewStart()
	*cmds = append(*cmds, startCmd)
	weatherCmd := command.NewWeather(sammy.config.GetString("configuration.weather"))
	*cmds = append(*cmds, weatherCmd)
	cnames = append(cnames, weatherCmd.Description())
	*cmds = append(*cmds, command.NewHelp(cnames))
	sammy.commands = cmds
}

func salute(salutations []string) (string) {
	return salutations[rand.Intn(len(salutations))]
}
