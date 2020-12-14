package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	tgBotAPI "gopkg.in/telegram-bot-api.v4"
)

func main() {
	port := os.Getenv("PORT")
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	bot, err := tgBotAPI.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal("creating bot failed: ", err)
	}

	log.Println("bot created")

	if _, err := bot.SetWebhook(tgBotAPI.NewWebhook(webHook)); err != nil {
		log.Fatalf("setting webHook: %v; error: %v", webHook, err)
	}

	log.Println("webHook set")

	updates := bot.ListenForWebhook("/")

	answers := []string{"Да", "Нет", "Возможно", "Затрудняюсь ответить", "Казалось бы"}
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(answers)

	for update := range updates {
		if _, err := bot.Send(tgBotAPI.NewMessage(update.Message.Chat.ID, answers[n])); err != nil {
			log.Print(err)
		}
	}
}
