package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	botpkg "summary_bot/bot"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка при загрузке .env файла.")
	}
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Println("TELEGRAM_BOT_TOKEN не установлен в .env файле.")
		os.Exit(1)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Println("Ошибка при создании бота", err)
		os.Exit(1)
	}
	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Ошибка при получении обновлений: %v", err)
	}
	for update := range updates {
		botpkg.HandleUpdate(bot, update)
	}
}
