package protectron

import (
	"log"
	"time"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Config struct {
	Token  string
	Repost time.Duration
	Link   time.Duration
	ProxyAuth
}

func Run(config Config) {
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
	}
}
