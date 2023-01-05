package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var nc map[int]string

func main() {
	wg := &sync.WaitGroup{}
	bot, u := mustBot(wg)
	updates := bot.GetUpdatesChan(u)

	play := false
	p := make([]int, 6)
	attempt := 0
	cn := map[string]int{"🟢": 0, "🔴": 1, "🟡": 2, "🔵": 3, "🟠": 4, "🟣": 5}
	nc = map[int]string{0: "🟢", 1: "🔴", 2: "🟡", 3: "🔵", 4: "🟠", 5: "🟣"}
	colors := ""
	for i := 0; i < len(nc); i++ {
		colors += nc[i]
	}
	for update := range updates {
		if update.Message != nil {
			msg := tgBotAPI.NewMessage(update.Message.Chat.ID, update.Message.Text)

			switch update.Message.Text {
			case "☑️":
				play = true
				attempt = 0
				rand.Seed(time.Now().UnixNano())
				p = rand.Perm(6)
				msg.Text = fmt.Sprintf("Введите комбинацию из 4-x разных цветов (%s):", colors)
				sendMsg(bot, msg)
			case "🔲":
				play = false
				msg.Text = "Стоп игра"
				sendMsg(bot, msg)
			default:
				if !play {
					continue
				}

				attempt++
				q := make([]int, 4)
				ans := strings.Split(msg.Text, "")
				for i := 0; i < 4; i++ {
					q[i] = cn[ans[i]]
					_, ok := cn[ans[i]]
					if !ok || len(ans) != 4 {
						msg.Text = fmt.Sprintf("Введите комбинацию из 4-x разных цветов (%s):", colors)
						sendMsg(bot, msg)
					}
				}

				msg.Text = checkAnswer(p, q, attempt)
				sendMsg(bot, msg)
			}
		} else if update.CallbackQuery != nil {
			if _, err := bot.Request(tgBotAPI.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)); err != nil {
				panic(err)
			}

			sendMsg(bot, tgBotAPI.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}
	}

	wg.Wait()
}

func mustBot(wg *sync.WaitGroup) (*tgBotAPI.BotAPI, tgBotAPI.UpdateConfig) {
	port := os.Getenv("PORT")
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}()

	bot, err := tgBotAPI.NewBotAPI(os.Getenv("tgToken"))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgBotAPI.NewUpdate(0)
	u.Timeout = 60

	return bot, u
}

func sendMsg(bot *tgBotAPI.BotAPI, msg tgBotAPI.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func checkAnswer(p, q []int, attempt int) string {
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

	text := fmt.Sprintf("%v", s)
	if win == 4 {
		line := ""
		for _, c := range p[:4] {
			line += nc[c]
		}
		text = fmt.Sprintf("Вы угадали c %d попытки!\nКомбинация: %s", attempt, line)
	}
	return text
}
