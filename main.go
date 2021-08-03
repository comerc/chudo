package main

import (
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	b *tb.Bot
)

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Ltime | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	var (
		// port      = os.Getenv("PORT")
		// publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		chats = os.Getenv("CHUDO_CHATS")  // you must add it to your config vars
		token = os.Getenv("CHUDO_SECRET") // you must add it to your config vars
	)
	log.Print(chats)
	// webhook := &tb.Webhook{
	// 	Listen:   ":" + port,
	// 	Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	// }
	pref := tb.Settings{
		// URL:    "https://api.bots.mn/telegram/",
		Token: token,
		// Poller: webhook,
		Poller:    &tb.LongPoller{Timeout: 10 * time.Minute},
		ParseMode: tb.ModeMarkdownV2,
	}
	{
		var err error
		b, err = tb.NewBot(pref)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Handle Ctrl+C
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		log.Print("Stop...")
		os.Exit(1)
	}()

	defer handlePanic()

	b.Handle(tb.OnCallback, func(*tb.Callback) {
		log.Print("OnCallback ", chats)
	})

	b.Handle(tb.OnQuery, func(q *tb.Query) {
		log.Print("OnQuery ", chats)
	})

	b.Handle(tb.OnPhoto, func(m *tb.Message) {
		log.Print("****")
		log.Print("Chat.ID ", m.Chat.ID)
		log.Print("****")

		// selector := &tb.ReplyMarkup{}
		// rows := make([]tb.Row, 0)
		// rows = append(rows, selector.Row(selector.Data("OK", "ok", "1234")))
		// selector.Inline(rows...)

		// _, err := b.Send(tb.ChatID(m.Chat.ID), `\.\.\.\.`, selector)
		// if err != nil {
		// 	log.Print(err)
		// }

		_, err := b.Forward(tb.ChatID(-1001168215421), m)
		if err != nil {
			log.Print(err)
		}

		// m2, err := b.Send(tb.ChatID(m.Chat.ID), "1111")
		// if err != nil {
		// 	log.Print(err)
		// }

		// _, err = b.Edit(m2, selector)
		// if err != nil {
		// 	log.Print(err)
		// }

		// _, err = b.Edit(m2, &tb.ReplyMarkup{})
		// if err != nil {
		// 	log.Print(err)
		// }

		// err := b.Delete(m)
		// if err != nil {
		// 	log.Print(err)
		// }
		// log.Print("OnText ", chats)

	})

	b.Handle("\fok", func(c *tb.Callback) {

		log.Print("ok ", c.Data)
	})

	// b.Handle("/start", func(q *tb.Message) {
	// 	log.Print("start ", chats)
	// })

	b.Start()
}

func handlePanic() {
	if err := recover(); err != nil {
		log.Printf("Panic...\n%s\n\n%s", err, debug.Stack())
		os.Exit(1)
	}
}
