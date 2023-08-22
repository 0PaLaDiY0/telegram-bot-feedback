package bot

import (
	"context"
	"fmt"
	"sync"
	l "telegram-bot-feedback/internal/pkg/logger"
	tg "telegram-bot-feedback/pkg/telegram-bot-api"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type App struct {
	Bot  *tg.Client
	DB   *gorm.DB
	Conf *viper.Viper
}

// Init initializes Telegram Bot
func Init(token, host string) (*tg.Client, error) {
	client, err := tg.NewWithHost(token, host)
	if err != nil {
		if err.Error() == "Not Found" {
			err = fmt.Errorf("incorrect token")
		}
		return nil, err
	}

	commandStart := tg.NewSetMyCommands(tg.BotCommand{Command: "/start", Description: "Starts chatting with the bot"})
	client.RequestOK(commandStart)

	return client, err
}

// RunFetcher handles Updates coming to the bot
func RunFetcher(ctx context.Context, wg *sync.WaitGroup, bot *tg.Client, db *gorm.DB, conf *viper.Viper) {
	defer wg.Done()
	app := App{Bot: bot, DB: db, Conf: conf}
	for {
		select {
		case <-ctx.Done():
			return
		default:
			updates := updates(bot, conf)
			for _, update := range updates {
				err := parseUpdate(&update, &app)
				if err != nil {
					l.Error(err)
					break
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// updates returns the slice of Update from the bot by offset
func updates(bot *tg.Client, conf *viper.Viper) []tg.Update {
	req := tg.NewUpdate(conf.GetInt("offset"))
	updates, err := bot.GetUpdates(req)
	if err != nil {
		l.Error(err)
		return nil
	}
	return updates
}
