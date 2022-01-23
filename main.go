package main

import (
	"log"
	"sync"

	"github.com/gadost/telescope/alert"
	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
)

var wg sync.WaitGroup

func main() {
	cfg, chains := conf.ConfLoad()
	log.Println(cfg)
	log.Println(chains)
	log.Println(chains[0])
	log.Println(conf.ConfdPath)
	log.Println(cfg.Chain["gravity"].Node[0].Role)
	log.Println(cfg.Chain["gravity"].Info.Mainnet)
	if conf.MainConfig.Telegram.Enabled {
		wg.Add(1)
		go alert.TelegramHandler()
	}
	watcher.ThreadsSplitter(cfg, chains)

}
