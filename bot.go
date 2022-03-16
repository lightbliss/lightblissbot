package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, u := mustBot()
	updates := bot.GetUpdatesChan(u)

	var err error
	play := false
	p := make([]int, 6)
	attempt := 0
	cn := map[string]int{"🟢": 0, "🔴": 1, "🟡": 2, "🔵": 3, "🟠": 4, "🟣": 5}
	nc := map[int]string{0: "🟢", 1: "🔴", 2: "🟡", 3: "🔵", 4: "🟠", 5: "🟣"}
	colors := ""
	for i := 0; i < len(nc); i++ {
		colors += nc[i]
	}
	var line string
	for update := range updates {
		if update.Message != nil {
			msg := tgBotAPI.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Text {
			case "play":
				play = true
				attempt = 0
				rand.Seed(time.Now().UnixNano())
				p = rand.Perm(6)
				msg.Text = fmt.Sprintf("Введите комбинацию разных 4-x цветов (%s):", colors)
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			case "stop":
				play = false
				msg.Text = "Стоп игра"
				if _, err = bot.Send(msg); err != nil {
					panic(err)
				}
			default:
				if play {
					attempt++
					q := make([]int, 4)
					ans := strings.Split(msg.Text, "")
					for i := 0; i < 4; i++ {
						q[i] = cn[ans[i]]
					}

					s := make([]string, 0)
					win := 0
					for i := 0; i < 4; i++ {
						if q[i] == p[i] {
							s = append(s, "⚫️")
							win++
						}
						for j := 0; j < 4; j++ {
							if i == j {
								continue
							}
							if q[i] == p[j] {
								s = append(s, "⚪️")
							}
						}
					}
					rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })

					msg.Text = fmt.Sprintf("%v", s)
					if win == 4 {
						line = ""
						for _, c := range p[:4] {
							line += nc[c]
						}
						msg.Text = fmt.Sprintf("Вы угадали c %d попытки! Комбинация: %s", attempt, line)
					}
					if _, err = bot.Send(msg); err != nil {
						panic(err)
					}
				}
			}
		} else if update.CallbackQuery != nil {
			callback := tgBotAPI.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			msg := tgBotAPI.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}

func mustBot() (*tgBotAPI.BotAPI, tgBotAPI.UpdateConfig) {
	port := os.Getenv("PORT")
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	bot, err := tgBotAPI.NewBotAPI(os.Getenv("tgToken"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 60

	return bot, u
}
