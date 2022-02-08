package status

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
	tele "gopkg.in/telebot.v3"
)

// StatusCollection collect status from Chain struct
func StatusCollection() string {
	var collection string
	var cu string
	var badRPC []string
	collection = collection + "*Status:*\n\n"

	for i := range watcher.Chains.Chain {
		for _, k := range watcher.Chains.Chain[i].Node {
			if k.Status.SyncInfo.LatestBlockHeight > 0 {
				collection += "*Net:* `" + k.Status.NodeInfo.Network + "`\n*Moniker:* `" + k.Status.NodeInfo.Moniker + "`\n"
				if k.Status.SyncInfo.CatchingUp {
					cu = "Yes"
				} else {
					cu = "No"
				}

				collection += "*CatchingUp:* `" + cu + "`\n"
				collection += "*Last known height:* `" + strconv.Itoa(int(k.Status.SyncInfo.LatestBlockHeight)) + "`\n"
				collection += "*Last known block time :* `" + k.Status.SyncInfo.LatestBlockTime.Format("2006-01-02 15:04:05") + "`\n"
				collection += "`_________________________`\n"
			} else {
				if k.MonitoringEnabled {
					badRPC = append(badRPC, k.RPC)
				}
			}
		}
	}

	if len(badRPC) > 0 {
		collection += "*🔴Unreachable RPCs:*\n`" + strings.Join(badRPC, "\n") + "`"
	}

	return collection
}

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
