package main

import (
	"fmt"
	"log"
	"os"
	"weatherbot/clients/openweather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	//

	// читает переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	owClient := openweather.New(os.Getenv("OPENWEATHERAPI_KEY"))
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, введите название города")
				bot.Send(msg)
				continue
			}

			coordinates, err := owClient.Coordinates(update.Message.Text)
			if err != nil {

				errorMsg := fmt.Sprintf("Не удалось найти город '%s'. Пожалуйста, проверьте название и попробуйте снова.", update.Message.Text)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}

			weather, err := owClient.Weather(coordinates.Lon, coordinates.Lat)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении погоды")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				fmt.Sprintf("Температура в %s: %.1f°C", update.Message.Text, weather.Temp))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
