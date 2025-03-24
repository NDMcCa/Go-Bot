package main

import (
	"fmt"

	"github.com/NDMcCa/Go-Bot/bot"
	"github.com/NDMcCa/Go-Bot/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})

	return
}
