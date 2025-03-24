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

var BotID string
var GoBot *discordgo.Session

func Start() {
	fmt.Println("Got key: ", config.WeatherKey)

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
	defer os.Exit(0)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Println("Goodbye!")

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}
	switch {
	case strings.Contains(m.Content, config.BotPrefix+"weather"):
		if strings.Contains(m.Content, "help") {
			_, _ = s.ChannelMessageSend(m.ChannelID, "To get the weather, type `"+config.BotPrefix+"weather city <city name> (-metric)` or `"+config.BotPrefix+"weather zip <zip code> (-metric)`")
			return
		}
		locType, locValue, units, err := parseWeatherCommand(m.Content)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		currentWeather := getCurrentWeather(locType, locValue, units)
		_, _ = s.ChannelMessageSendComplex(m.ChannelID, currentWeather)

	case m.Content == "<@"+BotID+"> ping" || m.Content == config.BotPrefix+"ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")

	case m.Content == "<@"+BotID+"> info" || m.Content == config.BotPrefix+"info":
		_, _ = s.ChannelMessageSend(m.ChannelID, "I am a Go Bot created by NDMcCa. My prefix is `"+config.BotPrefix+"`")

	case m.Content == "<@"+BotID+"> help" || m.Content == config.BotPrefix+"help":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Available Commands: "+"\n `"+config.BotPrefix+"weather`"+"\n `"+config.BotPrefix+"ping`")

	default:
		_, _ = s.ChannelMessageSend(m.ChannelID, "I'm sorry, I don't understand that command. Type `"+config.BotPrefix+"help` for a list of commands.")
	}

}
