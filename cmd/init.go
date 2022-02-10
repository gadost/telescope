/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gadost/telescope/alert"
	"github.com/gadost/telescope/conf"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init telescope",
	Long:  `Will create configurations directory and default <telescope.toml> config`,
	Run: func(cmd *cobra.Command, args []string) {
		Bootstrap()
		BootstrapSettings()
		BootstrapTelegram()
		Write()
	},
}
var (
	s        string
	r        = bufio.NewReader(os.Stdin)
	Telegram = false
	Discord  = false
	Twilio   = false
	Sms      = false
	Mail     = false
	Ghm      = false
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func Bootstrap() {
	if _, err := os.Stat(conf.UserHome + "/.telescope/conf.d"); !os.IsNotExist(err) {
		fmt.Println("config dir exist.")
		os.Exit(1)
	} else {
		if err := os.MkdirAll(conf.UserHome+"/.telescope/conf.d", os.ModePerm); err != nil {
			log.Fatal(err)
		}

	}

}

var TelescopeConfBasic = ``

func BootstrapSettings() {
	for {
		fmt.Fprint(os.Stderr, "Enable GitHub chain repository release monitoring? (Y/n)"+" ")
		s, _ = r.ReadString('\n')
		fmt.Print(s)
		if s == "Y\n" || s == "y\n" || s == "\n" {
			Ghm = true
			break
		} else if s == "N\n" || s == "n\n" {
			Ghm = false
			break
		}
	}

	TelescopeConfBasic += fmt.Sprintf(`
[settings]
github_release_monitor = %t
`, Ghm)

}

func BootstrapTelegram() bool {
	for {
		fmt.Fprint(os.Stderr, "Enable telegram notifications? (Y/n)"+" ")
		s, _ = r.ReadString('\n')
		if s == "Y\n" || s == "y\n" || s == "\n" {
			Telegram = true
			break
		} else if s == "N\n" || s == "n\n" {
			Telegram = false
			return Telegram
			break
		}
	}

	var token string
	var chatID string
	for {
		fmt.Fprint(os.Stderr, "Go to https://t.me/BotFather and create newbot. Enter bot token: ")
		token, _ = r.ReadString('\n')
		token = strings.TrimSuffix(token, "\n")
		fmt.Fprint(os.Stderr, "Enter chat id: ")
		chatID, _ = r.ReadString('\n')
		if chatID != "" {
			chatID = strings.TrimSuffix(chatID, "\n")
			fmt.Print("Sending ping.")
			err := alert.TelegramSendTest(token, chatID)
			if err != nil {
				fmt.Println(err)
			} else {
				var pass string
				fmt.Fprint(os.Stderr, " Ping received?(Y/n) ")
				pass, _ = r.ReadString('\n')
				if pass == "Y\n" || pass == "y\n" || pass == "\n" {
					break
				}
				fmt.Println("retry again.")
			}
		}

	}

	TelescopeConfBasic += fmt.Sprintf(`
[telegram]
enabled = %t
token = "%s"
chat_id = "%s"
`, Telegram, token, chatID)

	fmt.Print(TelescopeConfBasic)
	return Telegram
}

func BootstrapDiscord() {

}

func Write() {
	f, err := os.Create(conf.UserHome + "/.telescope/conf.d" + "/telescope.toml")

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data := []byte(TelescopeConfBasic)

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
