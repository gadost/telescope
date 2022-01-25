package alert

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gadost/telescope/conf"
	tele "gopkg.in/telebot.v3"
)

type Status struct {
}

func TelegramSend(s string, m string) {
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
	_, err = b.Send(to, "* "+s+"*: \n`"+m+"`", "MarkdownV2")
	if err != nil {
		fmt.Println(err)
	}

}
