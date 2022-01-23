package alert

import (
	"log"
	"strconv"
	"time"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
	tele "gopkg.in/telebot.v3"
)

type Status struct {
}

func TelegramHandler() {
	pref := tele.Settings{
		Token:  conf.MainConfig.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/status", func(c tele.Context) error {
		return c.Send(StatusCollection(), "HTML")
	})

	b.Start()
}

func StatusCollection() string {
	var collection string
	var cu string
	collection = collection + "<b>Status:</b>\n\n"
	for i := range watcher.Chains.Chain {
		for _, k := range watcher.Chains.Chain[i].Node {
			collection += "<b>Net:</b> " + k.Status.NodeInfo.Network + "\n<b>Moniker:</b> " + k.Status.NodeInfo.Moniker + "\n"
			if k.Status.SyncInfo.CatchingUp {
				cu = "Yes"
			} else {
				cu = "No"
			}

			collection += "<b>CatchingUp:</b> " + cu + "\n"
			collection += "<b>Last known height:</b> " + strconv.Itoa(int(k.Status.SyncInfo.LatestBlockHeight)) + "\n"
			collection += "<b>Last seen at:</b> " + k.Status.SyncInfo.LatestBlockTime.Format("2006-01-02 15:04:05") + "\n"
			collection += "_________________________\n"
		}
	}
	return collection
}
