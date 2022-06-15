package alert

import (
	"sync"

	"github.com/gadost/telescope/conf"
)

// wgAlert is wait group for alerts gorutines
var wgAlert sync.WaitGroup
var alertSystems = &conf.MainConfig
var at = "🔭 Telescope \\| "
var Importance = importance{
	Urgent:  at + "Urgent 🔴",
	Warning: at + "Warning ⚠️",
	Info:    at + "Info ℹ️",
	OK:      at + "OK ✅",
	GH:      at + "Repo Monitor 🔎",
}

// importance struct is struct for alert types
type importance struct {
	Urgent  string
	Warning string
	Info    string
	OK      string
	GH      string
}

// New creates new alert in gorutine
func New(i string, m string) {
	wgAlert.Add(1)
	go Alert(i, m)

	wgAlert.Wait()
}

// Alert send alerts to configured channels
func Alert(i string, m string) {
	defer wgAlert.Done()
	if alertSystems.Telegram.Enabled {
		TelegramSend(i, m)
	}
	if alertSystems.Discord.Enabled {
		DiscordSend(i, m)
	}
	/**
	if alertSystems.Mail.Enabled {

		}
		if alertSystems.Sms.Enabled {

		}
		if alertSystems.Twilio.Enabled {

		}
	**/
}
