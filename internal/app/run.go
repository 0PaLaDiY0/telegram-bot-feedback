package run

import (
	"context"
	"fmt"
	"os"
	"sync"
	tg "telegram-bot-feedback/internal/pkg/bot"
	"telegram-bot-feedback/internal/pkg/config"
	"telegram-bot-feedback/internal/pkg/console"
	"telegram-bot-feedback/internal/pkg/database"
	l "telegram-bot-feedback/internal/pkg/logger"
)

// Start starts bot
//
// Creates directories and configuration file
func Start() error {
	os.Mkdir("errors", 0755)

	conf, err := config.GetConfig()
	if err != nil {
		return l.Err(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	os.Mkdir("database", 0755)
	db, err := database.Init("database\\database.db")
	if err != nil {
		return l.Err(err)
	}

	if host := conf.GetString("host"); host == "" {
		fmt.Println("Enter the bot host in the format \"https://api.telegram.org/\" or \"-\" to use the standard:")
		fmt.Fscan(os.Stdin, &host)
		if host == "-" {
			conf.Set("host", "https://api.telegram.org/")
		} else {
			conf.Set("host", host)
		}
		if err := conf.WriteConfig(); err != nil {
			return l.Err(err)
		}
	}

	if token := conf.GetString("token"); token == "" {
		fmt.Println("Enter bot token:")
		fmt.Fscan(os.Stdin, &token)
		conf.Set("token", token)
		if err := conf.WriteConfig(); err != nil {
			return l.Err(err)
		}
	}

	client, err := tg.Init(conf.GetString("token"), conf.GetString("host"))
	if err != nil {
		return l.Err(err)
	}

	wg.Add(1)
	go tg.RunFetcher(ctx, &wg, client, db, conf)
	go console.Run(cancel, db)
	fmt.Println("Bot started")
	wg.Wait()
	return nil
}
