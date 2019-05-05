package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/klesogor/youtube-grabber/bot"
	"github.com/klesogor/youtube-grabber/grabber"
)

const videoUrl = "https://www.youtube.com/watch?v=cXEZu-uIdeI"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	grabber := grabber.NewHandler(10)
	bot.RunBot(os.Getenv("BOT_TOKEN"), grabber)
}
