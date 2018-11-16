package bot

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/proxy"
	"gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	token string
	proxy *SOCKS5ProxyConfig
}

type SOCKS5ProxyConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func (proxy *SOCKS5ProxyConfig) NoSecurity() bool {
	return proxy.Password == "" && proxy.Username == ""
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
		var auth *proxy.Auth
		if !bot.proxy.NoSecurity() {
			auth = &proxy.Auth{
				User:     bot.proxy.Username,
				Password: bot.proxy.Password,
			}
		}
		var dialer, errSOCKS5 = proxy.SOCKS5("tcp",
			fmt.Sprintf("%s:%d", bot.proxy.Host, bot.proxy.Port),
			auth, proxy.Direct)
		if errSOCKS5 != nil {
			return errSOCKS5
		}
		api, errAPI = tgbotapi.NewBotAPIWithClient(bot.token,
			&http.Client{
				Transport: &http.Transport{
					Dial: dialer.Dial,
				},
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
		case (update.Message != nil && update.Message.ForwardFromChat != nil):
			go func() {
				var userID = update.Message.From.ID
				var chatID = update.Message.Chat.ID
				var messageID = update.Message.MessageID
				log.Printf("[INFO] repost from @%s, engaging into battle", update.Message.ForwardFromChat.UserName)
				var resp, errKick = api.KickChatMember(tgbotapi.KickChatMemberConfig{
					ChatMemberConfig: tgbotapi.ChatMemberConfig{
						UserID: userID,
						ChatID: chatID,
					},
					UntilDate: -1,
				})
				if errKick != nil {
					log.Printf("[ERROR] %s", errKick)
					return
				}
				if !resp.Ok {
					log.Printf("[ERROR] %s", resp.Description)
					return
				}
				var respDelete, errDelete = api.DeleteMessage(tgbotapi.DeleteMessageConfig{
					ChatID:    chatID,
					MessageID: messageID,
				})
				if errDelete != nil {
					log.Printf("[ERROR] %s", errDelete)
					return
				}
				if !respDelete.Ok {
					log.Printf("[ERROR] %s", respDelete.Description)
					return
				}
			}()
		}
	}
	return nil
}
