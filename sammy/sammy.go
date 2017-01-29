package sammy

import (
	"math/rand"
	"log"

	"github.com/spf13/viper"
	"github.com/MarcosSegovia/sammy-the-bot/command"
)

type Request string
type Response string

func (resp Response) String() string {
	return string(resp)
}

type Sammy struct {
	brain    *viper.Viper
	commands map[string]interface{}
}

func NewSammySpeaker(brain *viper.Viper) *Sammy {
	s := new(Sammy)
	s.brain = brain
	s.load()
	return s
}

func (sammy *Sammy) Process(req Request) Response {
	resp := Response("I do not know what to tell you.")
	if "Hi" == req {
		salutations := sammy.brain.GetStringSlice("welcome.salutations")
		resp = Response(salute(salutations))
	}
	for i, v := range sammy.commands {
		if "start" == i || "help" == i {
			cmd := v.(command.Cmd)
			if string(cmd.Exec) == string(req) {
				resp = sammy.ProcessCmd(&cmd)
			}
		}
	}
	log.Printf("I'm responding: %v", resp)
	return resp
}

func (sammy *Sammy) ProcessCmd(cmd command.Command) Response {
	buffer := cmd.Evaluate()
	return Response(buffer.String())
}

func (sammy *Sammy) load() {
	var commands = make(map[string]interface{}, 2)
	commands["start"] = command.Cmd{Tag:"start", Exec: "/start"}
	commands["help"] = command.Cmd{Tag:"help", Exec: "/help"}
	sammy.commands = commands
}

func salute(salutations []string) (string) {
	return salutations[rand.Intn(len(salutations))]
}
