package main

import (
	"log"
	"sync"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
)

var wgMain sync.WaitGroup

func main() {
	cfg, chains := conf.ConfLoad()
	if conf.MainConfig.Telegram.Enabled {
		log.Println("Telegram command handler started.")
		wgMain.Add(1)
		go watcher.TelegramHandler()
	}
	if conf.MainConfig.Settings.GithubReleaseMonitor {
		log.Println("Github repositories monitor started.")
		wgMain.Add(1)
		go watcher.CheckNewRealeases()
	}
	wgMain.Add(1)
	go watcher.BlockProducingParticipation(cfg, chains)
	log.Println("Alert system started for chains:", chains)
	watcher.ThreadsSplitter(cfg, chains)

	wgMain.Wait()
}
