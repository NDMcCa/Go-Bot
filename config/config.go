package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	Token      string
	BotPrefix  string
	WeatherKey string

	config *Config
)

type Config struct {
	Token      string `json:"token"`
	BotPrefix  string `json:"BotPrefix"`
	WeatherKey string `json:"weatherKey"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := os.ReadFile("./config.json")

	if err != nil {
		return err
	}

	fmt.Println("Config file opened successfully")
	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println("Error Unmarshalling config file")
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix
	WeatherKey = config.WeatherKey

	return nil
}
