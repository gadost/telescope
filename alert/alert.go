package alert

import (
	"sync"

	"github.com/gadost/telescope/conf"
)

var wg sync.WaitGroup
var alertSystems = &conf.MainConfig
var at = "🔭 Telescope \\| "
var Importance = importance{
	Urgent:  at + "Urgent🔴",
	Warning: at + "Warning ⚠️",
	Info:    at + "Info ℹ️",
	OK:      at + "OK ✅",
}

type importance struct {
	Urgent  string
	Warning string
	Info    string
	OK      string
}

func New(i string, m string) {
	wg.Add(1)
	go Alert(i, m)
}

func Alert(i string, m string) {
	defer wg.Done()
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
