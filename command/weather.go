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
	defer resp.Body.Close()
	check(err, "could not get an appropriate response: %v")

	body, err := ioutil.ReadAll(resp.Body)
	wresp := Response{}
	err = json.Unmarshal(body, &wresp)
	check(err, "could not parse from json: %v")
	buffer.Reset()

	buffer.WriteString("Seems that we will have ")
	buffer.WriteString(wresp.Conditions[0]["main"].(string) + " ")
	switch  wresp.Conditions[0]["id"].(float64) {

	//Thunder
	case 200, 201, 202, 210, 211, 212, 221, 230, 231, 232:
		buffer.Write([]byte{226, 155, 136})

	//Drizzle
	case 300, 301, 302, 310, 311, 312, 313, 314, 321:
		buffer.Write([]byte{226, 155, 136})

	//Rain
	case 500, 501, 502, 503, 504, 511, 520, 521, 522, 531:
		buffer.Write([]byte{240, 159, 140, 167})

	//Snow
	case 600, 601, 602, 611, 612, 615, 616, 620, 621, 622:
		buffer.Write([]byte{226, 157, 132})

	//Fog
	case 701, 711, 721, 731, 741, 751, 761, 762, 771, 781:
		buffer.Write([]byte{240, 159, 140, 171})

	//Clear
	case 800:
		buffer.Write([]byte{226, 152, 128})

	//Few clouds
	case 801:
		buffer.Write([]byte{240, 159, 140, 164})

	//Clouds
	case 802, 803, 804:
		buffer.Write([]byte{226, 155, 133})

	}
	buffer.WriteString(" for today here in ")
	buffer.WriteString(wresp.City)
	buffer.WriteString(", with a temperature of ")
	celsius := tempconv.KelvinToCelcius(tempconv.Kelvin(wresp.Main["temp"]))
	buffer.WriteString(celsius.String())
	return buffer
}

func check(err error, msg string) {
	if err != nil {
		log.Printf(msg, err)
	}
}
