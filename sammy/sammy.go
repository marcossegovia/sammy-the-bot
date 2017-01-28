package sammy

import (
	"math/rand"
	"log"

	"github.com/spf13/viper"
)

type Request string
type Response string

func (resp Response) String() string {
	return string(resp)
}

type Sammy struct {
	brain *viper.Viper
}

func (sammy *Sammy) Process(req Request) Response {
	resp := Response("I do not know what to tell you.")
	if "Hi" == req {
		salutations := sammy.brain.GetStringSlice("welcome.salutations")
		resp = Response(salute(salutations))
	}
	log.Printf("I'm responding: %v", resp)

	return resp
}

func NewSammySpeaker(brain *viper.Viper) *Sammy {
	s := new(Sammy)
	s.brain = brain
	return s
}

func salute(salutations []string) (string) {
	return salutations[rand.Intn(len(salutations))]
}
