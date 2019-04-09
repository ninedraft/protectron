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

func (config Config) getMaxQuarantine() time.Duration {
	var qd = []time.Duration{
		config.Repost,
		config.Link,
	}
	var max time.Duration
	for _, q := range qd {
		if q > max {
			max = q
		}
	}
	return max
}

func Run(config Config) {
	var reg = newRegistry(config.getMaxQuarantine())
	defer reg.stopVacuum()
	bot, errNewBot := tgbotapi.NewBotAPI(config.Token)
	if errNewBot != nil {
		log.Fatal(errNewBot)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, errGetUpdates := bot.GetUpdatesChan(u)
	if errGetUpdates != nil {
		log.Fatal(errGetUpdates)
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		go func(msg *tgbotapi.Message) {
			defer func() {
				var panicMsg = recover()
				if panicMsg != nil {
					log.Printf("recovered: %v", panicMsg)
				}
			}()
			var userID = msg.From.ID
			var repostedTooEarly = msg.ForwardFromChat != nil && reg.userIsTooYoung(config.Repost, userID)
			var postedLinkToEarly = msgHasLinks(msg) && reg.userIsTooYoung(config.Link, userID)
			switch {
			case repostedTooEarly, postedLinkToEarly:
				log.Printf("kicking user %q:%d from chat", msg.From.UserName, userID)
				var _, errKickChatMember = bot.KickChatMember(tgbotapi.KickChatMemberConfig{
					ChatMemberConfig: tgbotapi.ChatMemberConfig{
						ChatID: msg.Chat.ID,
						UserID: userID,
					},
				})
				if errKickChatMember != nil {
					log.Printf("unable to ban user @%s: %v", msg.From.UserName, errKickChatMember)
				}
			case msg.NewChatMembers != nil:
				for _, newUser := range *msg.NewChatMembers {
					log.Printf("user %q:%d joined chat", newUser.UserName, userID)
					reg.addUser(newUser.ID)
				}
			}
		}(update.Message)
	}
}

func msgHasLinks(msg *tgbotapi.Message) bool {
	if msg.Entities == nil {
		return false
	}
	for _, entity := range *msg.Entities {
		var hasLink = entity.URL != ""
		var hasBotUsername = entity.User != nil && entity.User.IsBot
		switch {
		case hasLink, hasBotUsername:
			return true
		}
	}
	return false
}
