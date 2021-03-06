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

	tgToken := os.Getenv("tgToken")
	bot, err := tgBotAPI.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal("creating bot failed: ", err)
	}

	log.Println("bot created")

	webHook := os.Getenv("webHook")
	if _, err := bot.SetWebhook(tgBotAPI.NewWebhook(webHook)); err != nil {
		log.Fatalf("setting webHook: %v; error: %v", webHook, err)
	}

	log.Println("webHook set")

	updates := bot.ListenForWebhook("/")

	answers := []string{"Да", "Определенно", "Нет", "Никогда", "Возможно", "Казалось бы", "Когда-нибудь"}

	for update := range updates {
		rand.Seed(int64(time.Now().Nanosecond()))
		n := rand.Intn(len(answers))
		log.Printf("answer is %v", answers[n])
		if _, err := bot.Send(tgBotAPI.NewMessage(update.Message.Chat.ID, answers[n])); err != nil {
			log.Print(err)
		}
	}
}
