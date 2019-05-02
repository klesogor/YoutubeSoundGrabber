package bot

import (
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/klesogor/youtube-grabber/grabber"
)

const processMessage = "Your request have been added to process queue. Please be patient, while we converting your audio..."

func RunBot(token string, youtubeGrabber grabber.YoutubeGrabber) {
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

		messageHandler, fileHandler := createMesageHandler(chatId, messageId, bot), createFileHandler(chatId, messageId, bot)
		handlers := grabber.Handlers{MessageHandler: messageHandler, FileHandler: fileHandler}

		youtubeGrabber.Handle(update.Message.Text, handlers)
	}
}

func createMesageHandler(chatId int64, messageId int, bot *tgbotapi.BotAPI) grabber.MessageHandler {
	return func(mes grabber.ResponseMessage) {
		msg := tgbotapi.NewMessage(chatId, mes.Err.Error())
		msg.ReplyToMessageID = messageId

		bot.Send(msg)
	}
}

func createFileHandler(chatId int64, messageId int, bot *tgbotapi.BotAPI) grabber.FileMessageHandler {
	return func(mes grabber.ResponseFileMessage) {
		msg := tgbotapi.NewAudioUpload(chatId, mes.FilePath)
		msg.ReplyToMessageID = messageId

		bot.Send(msg)
	}
}
