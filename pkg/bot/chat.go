package bot

import (
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func defend(ctx botContext) error {
	var senderID = ctx.Msg.From.ID
	// var senderUsername = ctx.Msg.From.UserName
	var msgID = ctx.Msg.MessageID
	var chatID = ctx.Msg.Chat.ID
	var _, errDelete = ctx.Bot.DeleteMessage(tgbotapi.DeleteMessageConfig{
		MessageID: msgID,
		ChatID:    chatID,
	})
	if errDelete != nil {
		return errDelete
	}
	var _, errKick = ctx.Bot.KickChatMember(tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: senderID,
		},
		UntilDate: -1,
	})
	return errKick
}

func processMessage(ctx botContext) {
	defer func() {
		if err := recover(); err != nil {
			logErr("recovered: %v", err)
		}
	}()
	if isForward(ctx.Msg, ctx.ChatWhitelist) ||
		(containsMedia(ctx.Msg) && constainsSuspiciousURL(ctx.Msg, ctx.HostWhitelist)) {
		if err := defend(ctx); err != nil {
			logErr("unable to ban user: %v", err)
		}
	}
}
