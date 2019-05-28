package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/klesogor/youtube-grabber/bot"
	"github.com/klesogor/youtube-grabber/internals/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cache := telegram.NewMongoCache(os.Getenv("MONGO_CONNECTION"))
	fmt.Println("saved test data to cache")
	bot.RunBot(os.Getenv("BOT_TOKEN"), cache)
}
