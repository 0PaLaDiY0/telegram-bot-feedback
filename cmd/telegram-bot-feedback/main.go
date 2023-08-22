package main

import (
	"fmt"
	bot "telegram-bot-feedback/internal/app"
	l "telegram-bot-feedback/internal/pkg/logger"
)

// Starts the bot
func main() {
	err := bot.Start()
	if err != nil {
		l.Fatal(err)
		fmt.Println("Launch error")
		fmt.Print("Press Enter to complete")
		fmt.Scanln()
		return
	}
}
