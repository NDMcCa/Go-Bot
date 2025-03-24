package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

const URL string = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Name string `json:"name"`
}

func getCurrentWeatherZIP(message string, units string) *discordgo.MessageSend {
	r, _ := regexp.Compile(`!zip (\d{5})`)
	zip := r.FindStringSubmatch(message)[1]

	if zip == "" {
		return &discordgo.MessageSend{
			Content: "Invalid zip code",
		}
	}

	weatherURL := fmt.Sprintf("%szip=%s&appid=%s&units=%s", URL, zip, WeatherToken, units)

	client := http.Client{
		Timeout: time.Second * 2,
	}

	response, err := client.Get(weatherURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Error getting weather data",
		}
	}

	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var data WeatherData
	json.Unmarshal([]byte(body), &data)

	location := data.Name
	conditions := data.Weather[0].Description
	temp := strconv.FormatFloat(data.Main.Temp, 'f', 2, 64)
	pressure := strconv.Itoa(data.Main.Pressure)
	humidity := strconv.Itoa(data.Main.Humidity)
	windSpeed := strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64)
	windDirection := strconv.Itoa(data.Wind.Deg)

	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title:       "Current Weather",
			Description: "Weather for " + location,
			Fields: []*discordgo.MessageEmbedField{{
				Name:   "Conditions",
				Value:  conditions,
				Inline: true,
			}, {
				Name:   "Temperature",
				Value:  temp + "°F",
				Inline: true,
			}, {
				Name:   "Pressure",
				Value:  pressure + "hPa",
				Inline: true,
			}, {
				Name:   "Humidity",
				Value:  humidity + "%",
				Inline: true,
			}, {
				Name:   "Wind Speed",
				Value:  windSpeed + "mph",
				Inline: true,
			}, {
				Name:   "Wind Direction",
				Value:  windDirection + "°",
				Inline: true,
			}},
		}},
	}
	return embed
}
