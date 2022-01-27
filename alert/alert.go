package alert

import (
	"sync"

	"github.com/gadost/telescope/conf"
)

var wgAlert sync.WaitGroup
var alertSystems = &conf.MainConfig
var at = "🔭 Telescope \\| "
var Importance = importance{
	Urgent:  at + "Urgent🔴",
	Warning: at + "Warning ⚠️",
	Info:    at + "Info ℹ️",
	OK:      at + "OK ✅",
	GH:      at + "Repo Monitor 🔎",
}

type importance struct {
	Urgent  string
	Warning string
	Info    string
	OK      string
	GH      string
}

func New(i string, m string) {
	wgAlert.Add(1)
	go Alert(i, m)

	wgAlert.Wait()
}

func Alert(i string, m string) {
	defer wgAlert.Done()
	if alertSystems.Telegram.Enabled {
		TelegramSend(i, m)
	}
	/**	if alertSystems.Discord.Enabled {

		}
		if alertSystems.Mail.Enabled {

		}
		if alertSystems.Sms.Enabled {

		}
		if alertSystems.Twilio.Enabled {

		}
	**/
}
