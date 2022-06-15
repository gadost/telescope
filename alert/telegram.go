package alert

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gadost/telescope/conf"
	tele "gopkg.in/telebot.v3"
)

// TelegramSend send message to configured telegram channel
func (a *Alert) TelegramSend() {
	var pref = tele.Settings{
		Token:  conf.MainConfig.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	var chatID, _ = strconv.ParseInt(conf.MainConfig.Telegram.ChatID, 10, 64)
	var to = &tele.Chat{ID: chatID}
	_, err = b.Send(to, "* "+a.Importance+"*: \n`"+a.Message+"`", "MarkdownV2")
	if err != nil {
		fmt.Println(err)
	}
}

func TelegramSendTest(t string, c string) error {
	var pref = tele.Settings{
		Token:  t,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		return err
	}
	var chatID, _ = strconv.ParseInt(c, 10, 64)
	var to = &tele.Chat{ID: chatID}
	_, err = b.Send(to, "Pong")
	return err
}
