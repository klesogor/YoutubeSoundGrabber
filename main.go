package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/klesogor/youtube-grabber/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot.RunBot(os.Getenv("BOT_TOKEN"))
}
