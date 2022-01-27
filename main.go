package main

import (
	"log"
	"sync"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
)

var wg sync.WaitGroup

func main() {
	cfg, chains := conf.ConfLoad()
	if conf.MainConfig.Telegram.Enabled {
		log.Println("Telegram command handler started.")
		wg.Add(1)
		go watcher.TelegramHandler()
	}
	if conf.MainConfig.Settings.GithubReleaseMonitor {
		log.Println("Github repositories monitor started.")
		wg.Add(1)
		go watcher.CheckNewRealeases()
	}
	log.Println("Alert system started for chains:", chains)
	watcher.ThreadsSplitter(cfg, chains)

}
