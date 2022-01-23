package main

import (
	"sync"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
)

var wg sync.WaitGroup

func main() {
	cfg, chains := conf.ConfLoad()
	if conf.MainConfig.Telegram.Enabled {
		wg.Add(1)
		go watcher.TelegramHandler()
	}
	watcher.ThreadsSplitter(cfg, chains)

}
