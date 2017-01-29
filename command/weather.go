package command

import (
	"io/ioutil"
	"net/http"
	"log"
	"bytes"
	"encoding/json"

	"github.com/danbondd/temperature/tempconv"
)

type Weather struct {
	Cmd   *Cmd
	Token string
}

type Response struct {
	Coord      map[string]float64 `json:"coord"`
	Conditions []map[string]interface{} `json:"weather"`
	Main       map[string]float64 `json:"main"`
	City       string `json:"name"`
}

func NewWeatherCommand(token string) (*Weather) {
	wCmd := new(Weather)
	wCmd.Cmd = NewCommand("weather", "/weather")
	wCmd.Token = token
	return wCmd
}

func (w *Weather) Evaluate() bytes.Buffer {
	var buffer bytes.Buffer
	buffer.WriteString("http://api.openweathermap.org/data/2.5/weather?id=6356055&appid=")
	buffer.WriteString(w.Token)
	resp, err := http.Get(buffer.String())
	check(err, "could not get an appropriate response: %v")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	wresp := Response{}
	err = json.Unmarshal(body, &wresp)
	check(err, "could not parse from json: %v")
	buffer.Reset()
	buffer.WriteString("Seems that we will have ")
	buffer.WriteString(wresp.Conditions[0]["main"].(string))
	buffer.WriteString(" for today here in ")
	buffer.WriteString(wresp.City)
	buffer.WriteString(", with a temperature of ")
	celsius := tempconv.KelvinToCelcius(tempconv.Kelvin(wresp.Main["temp"]))
	buffer.WriteString(celsius.String())
	buffer.WriteString(" ")
	buffer.WriteString("http://openweathermap.org/img/w/")
	buffer.WriteString(wresp.Conditions[0]["icon"].(string))
	buffer.WriteString(".png")
	return buffer
}

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
