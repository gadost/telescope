/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"
	"sync"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/status"
	"github.com/gadost/telescope/watcher"
	"github.com/spf13/cobra"
)

var wgMain sync.WaitGroup

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Telescope service",
	Long: `Start Telescope service
	By default telescope using ~/.telescope/conf.d as configurations directory. 
	You can specify by passsing --conf flag`,
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}
var c string

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVar(&c, "conf", conf.UserHome+"/.telescope/conf.d", "Configurations directory")

}

func Start() {
	cfg, chains := conf.ConfLoad(c)
	if conf.MainConfig.Telegram.Enabled {
		log.Println("Telegram commands handler bot started.")
		wgMain.Add(1)
		go status.TelegramHandler()
	}
	if conf.MainConfig.Discord.Enabled {
		log.Println("Discord commands handler bot started.")
		wgMain.Add(1)
		go status.DiscordHandler()
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
