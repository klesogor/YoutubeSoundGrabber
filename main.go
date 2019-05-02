package main

import (
	"fmt"

	"github.com/klesogor/youtube-grabber/bot"
	"github.com/klesogor/youtube-grabber/grabber"
)

func main() {
	fmt.Println("Starting bot...")
	bot.RunBot("", grabber.NewHandler(5))
}
