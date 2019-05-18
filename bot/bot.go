package bot

import (
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

const processMessage = "Your request have been added to process queue. Please be patient, while we converting your audio..."

type TelegramContext struct {
	chatId    int64
	messageId int
}

func RunBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		chatId, messageId := update.Message.Chat.ID, update.Message.MessageID
		msg := tgbotapi.NewMessage(chatId, processMessage)
		msg.ReplyToMessageID = messageId
		bot.Send(msg)
	}
}
