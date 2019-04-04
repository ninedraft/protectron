package bot

import (
	"log"
	"net/http"
	"os"

	"github.com/ninedraft/protectron/pkg/proxy"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	token string
	proxy *SOCKS5ProxyConfig
}

func New(token string, proxy *SOCKS5ProxyConfig) Bot {
	return Bot{
		token: token,
		proxy: proxy,
	}
}

func (bot Bot) Run() error {
	log.Printf("connecting to Bot API")
	const timeout = 10 // seconds
	var api *tgbotapi.BotAPI
	var errAPI error

	if bot.proxy == nil {
		api, errAPI = tgbotapi.NewBotAPI(bot.token)
		if errAPI != nil {
			return errAPI
		}
	} else {
		log.Printf("using proxy")
		var client, errConnectProxy = proxy.SOCKS5(bot.proxy.Address(), proxy.Auth{
			User:     bot.proxy.Username,
			Password: bot.proxy.Password,
		})
		if errConnectProxy != nil {
			logErr("unable to use proxy: %v", errConnectProxy)
			os.Exit(1)
		}
		api, errAPI = tgbotapi.NewBotAPIWithClient(bot.token,
			&http.Client{
				Transport: client.Transport,
			})
		if errAPI != nil {
			return errAPI
		}
	}

	log.Printf("getting updates")
	var updates, errUpdates = api.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 10,
	})
	if errUpdates != nil {
		return errUpdates
	}
	log.Printf("listening for updates")
	for update := range updates {
		switch {
		case update.Message != nil:
			go processMessage(botContext{
				Bot:           api,
				Msg:           update.Message,
				HostWhitelist: strSet([]string{}),
				ChatWhitelist: func(id int64) bool { return true },
			})
		default:
			continue
		}
	}
	return nil
}

type botContext struct {
	Bot           *tgbotapi.BotAPI
	Msg           *tgbotapi.Message
	ChatWhitelist idFilter
	HostWhitelist strFilter
}
