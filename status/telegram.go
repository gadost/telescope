package status

import (
	"log"
	"time"

	"github.com/gadost/telescope/conf"
	tele "gopkg.in/telebot.v3"
)

// TelegramHandler start telegram command handler
func TelegramHandler() {
	var pref = tele.Settings{
		Token:  conf.MainConfig.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/status", func(c tele.Context) error {
		return c.Send(StatusCollection(), "MarkdownV2")
	})

	b.Start()
}
