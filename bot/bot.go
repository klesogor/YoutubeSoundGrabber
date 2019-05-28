package bot

import (
	"log"

	"github.com/klesogor/youtube-grabber/internals/telegram"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/klesogor/youtube-grabber/internals"
	"github.com/klesogor/youtube-grabber/internals/youtube"
)

const processMessage = "Hi! I'm youtube converter, just send me youtube url, ad I will extract audio from it."

var converter internals.FFMPEGConverter

type TelegramContext struct {
	chatId    int64
	messageId int
}

func RunBot(token string, cache telegram.TelegramAudioCache) {
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
		case "/start":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, processMessage))
			break
		default:
			go processVideo(bot, update.Message, cache)
		}
	}
}

func processVideo(bot *tgbotapi.BotAPI, message *tgbotapi.Message, cache telegram.TelegramAudioCache) {
	conf, err := youtube.GetPlayerConfig(message.Text)
	if err != nil {
		reportError(bot, err, message)
		return
	}
	if audioId, err := cache.TryGetAudioId(conf.Args.VideoID); err == nil {
		bot.Send(tgbotapi.NewAudioShare(message.Chat.ID, audioId))
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
	bytes := tgbotapi.FileBytes{Name: conf.Args.Title, Bytes: converted}
	res, err := bot.Send(tgbotapi.NewAudioUpload(message.Chat.ID, bytes))
	if err == nil {
		cache.SaveAudioIdToCache(conf.Args.VideoID, res.Audio.FileID)
	}
}
func reportError(bot *tgbotapi.BotAPI, err error, message *tgbotapi.Message) {
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, err.Error()))
}
