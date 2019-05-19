package bot

import (
	"log"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/klesogor/youtube-grabber/internals"
	"github.com/klesogor/youtube-grabber/internals/youtube"
)

const processMessage = "Your request have been added to process queue. Please be patient, while we converting your audio..."

var converter internals.FFMPEGConverter

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
		switch update.Message.Text {

		}

		chatId, messageId := update.Message.Chat.ID, update.Message.MessageID
		msg := tgbotapi.NewMessage(chatId, processMessage)
		msg.ReplyToMessageID = messageId
		bot.Send(msg)
	}
}

func processVideo(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	conf, err := youtube.GetPlayerConfig(message.Text)
	if err != nil {
		reportError(bot, err, message)
		return
	}
	audio, err := conf.DownloadAudio()
	if err != nil {
		reportError(bot, err, message)
		return
	}
	converted, err := converter.Convert(audio, internals.ConvertingSettings{PreserveVideo: false, TargetFormat: internals.MP3})
	if err != nil {
		reportError(bot, err, message)
		return
	}

	bot.Send(tgbotapi.NewAudioUpload(message.Chat.ID, converted))
}
func reportError(bot *tgbotapi.BotAPI, err error, message *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
}
