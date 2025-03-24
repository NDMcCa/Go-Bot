package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/NDMcCa/Go-Bot/config"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID        string
	WeatherToken string = "b6ee0cf951e243fe3a840e606d01ba54"
)

var GoBot *discordgo.Session

func Start() {
	fmt.Println("Got key: ", WeatherToken)

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	defer goBot.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	switch {
	case strings.Contains(m.Content, "weather"):
		_, _ = s.ChannelMessageSend(m.ChannelID, "Weather is not implemented yet")

	case strings.Contains(m.Content, config.BotPrefix+"!zip"):
		if strings.Contains(m.Content, "metric") {
			currentWeather := getCurrentWeatherZIP(m.Content, "metric")
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, currentWeather)
		} else {
			currentWeather := getCurrentWeatherZIP(m.Content, "imperial")
			_, _ = s.ChannelMessageSendComplex(m.ChannelID, currentWeather)
		}

	case m.Content == "<@"+BotID+"> ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")

	case m.Content == "<@"+BotID+"> help":
		_, _ = s.ChannelMessageSend(m.ChannelID, "I am a Go Bot created by NDMcCa. My prefix is `"+config.BotPrefix+"`")

	case m.Content == config.BotPrefix+"ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")

	}
}
