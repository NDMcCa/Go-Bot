package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/NDMcCa/Go-Bot/config"
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

func parseWeatherCommand(message string) (string, string, string, error) {
	if strings.Contains(message, "zip") {
		r, _ := regexp.Compile(`!weather zip (\d{5})( -metric)?`)
		match := r.FindStringSubmatch(message)
		if len(match) > 0 {
			zip := match[1]
			units := "imperial"
			if len(match) > 2 && strings.TrimSpace(match[2]) == "-metric" {
				units = "metric"
			}
			return "zip", zip, units, nil
		}
		return "", "", "", fmt.Errorf("invalid zip format")
	} else if strings.Contains(message, "city") {
		r, _ := regexp.Compile(`!weather city (.+?)( -metric)?$`)
		match := r.FindStringSubmatch(message)
		if len(match) > 0 {
			city := match[1]
			units := "imperial"
			if len(match) > 2 && strings.TrimSpace(match[2]) == "-metric" {
				units = "metric"
			}
			return "city", city, units, nil
		}
		return "", "", "", fmt.Errorf("invalid city format")
	}
	return "", "", "", fmt.Errorf("invalid command")
}

func getCurrentWeather(locType string, locValue string, units string) *discordgo.MessageSend {
	var weatherURL string
	if locType == "zip" {
		weatherURL = fmt.Sprintf("%szip=%s&appid=%s&units=%s", URL, locValue, config.WeatherKey, units)
	} else {
		weatherURL = fmt.Sprintf("%sq=%s&appid=%s&units=%s", URL, locValue, config.WeatherKey, units)
	}

	fmt.Println("Fetching weather from URL: " + weatherURL)

	client := http.Client{
		Timeout: time.Second * 2,
	}

	response, err := client.Get(weatherURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Error getting weather data",
		}
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	var data WeatherData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Error parsing weather data",
		}
	}

	if len(data.Weather) == 0 {
		return &discordgo.MessageSend{
			Content: "No weather data available; the location may not exist",
		}
	}

	speed := "mph"
	if units == "metric" {
		speed = "km/h"
	}

	temp_unit := "°F"
	if units == "metric" {
		temp_unit = "°C"
	}

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
				Value:  temp + temp_unit,
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
				Value:  windSpeed + speed,
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
