package main

import (
	"log"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/watcher"
)

func main() {
	cfg, chains := conf.ConfLoad()
	log.Println(cfg)
	log.Println(chains)
	log.Println(chains[0])
	log.Println(cfg.Chain["gravity"].Node[0].Role)
	log.Println(cfg.Chain["gravity"].Info.Mainnet)
	watcher.ThreadsSplitter(cfg, chains)
	log.Println(chains[0])

}
