package bot

import (
	"net/url"

	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func urlsFromMessage(msg *tgbotapi.Message) []*url.URL {
	if msg.Entities == nil {
		return []*url.URL{}
	}
	var urls = make([]*url.URL, 0, len(*msg.Entities))
	for _, entity := range *msg.Entities {
		if entity.URL != "" {
			var innerURL, _ = entity.ParseURL()
			urls = append(urls, innerURL)
		}
	}
	return urls
}

func constainsSuspiciousURL(msg *tgbotapi.Message, hostWhitelist strFilter) bool {
	if msg.Entities == nil {
		return false
	}
	for _, entity := range *msg.Entities {
		if entity.URL != "" {
			var innerURL, _ = entity.ParseURL()
			return !hostWhitelist(innerURL.Host)
		}
	}
	return false
}

func containsMedia(msg *tgbotapi.Message) bool {
	return msg.Photo != nil || msg.Video != nil
}

func isForward(msg *tgbotapi.Message, chatWhiteist idFilter) bool {
	return msg.ForwardFromChat != nil && !chatWhiteist(msg.ForwardFromChat.ID)
}
