package bot

import (
	"fmt"

	"github.com/NDMcCa/Go-Bot/config"

	"github.com/bwmarrin/discordgo"
)

var BotID string
var GoBot *discordgo.Session

func Start() {
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

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if m.Content == "<@"+BotID+"> ping" || m.Content == "<@"+BotID+">" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == config.BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
