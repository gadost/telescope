package alert

import (
	"fmt"
	"sync"
	"time"

	"github.com/gadost/telescope/conf"
)

const (
	at      = "🔭 Telescope \\| "
	urgent  = at + "Urgent 🔴"
	warning = at + "Warning ⚠️"
	info    = at + "Info ℹ️"
	ok      = at + "OK ✅"
	gh      = at + "Repo Monitor 🔎"
)

// wgAlert is wait group for alerts gorutines
var wgAlert sync.WaitGroup
var alertSystems = &conf.MainConfig
var Importance = importance{
	Urgent:  urgent,
	Warning: warning,
	Info:    info,
	OK:      ok,
	GH:      gh,
}

// importance struct is struct for alert types
type importance struct {
	Urgent  string
	Warning string
	Info    string
	OK      string
	GH      string
}

type Alert struct {
	Importance string
	Message    string
}

func New(i, m string) *Alert {
	return &Alert{Importance: i, Message: m}
}

// New creates new alert in gorutine
func (a *Alert) Send() {
	wgAlert.Add(1)
	go func(a *Alert) {
		defer wgAlert.Done()
		if alertSystems.Telegram.Enabled {
			a.TelegramSend()
		}
		if alertSystems.Discord.Enabled {
			a.DiscordSend()
		}
	}(a)

	wgAlert.Wait()
}

func NewAlertBlockMissed(moniker string, network string, missed int) *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Validator '%s' , \nNetwork: %s \nSignature missed  in last %v blocks in a row",
		moniker,
		network,
		missed,
	)
	return New(i, m)
}

func NewAlertVotingPower(moniker string, network string, diff string) *Alert {
	i := Importance.Info
	m := fmt.Sprintf("Voting Power of '%s' , \nNetwork: %s \n%s", moniker, network, diff)
	return New(i, m)
}

func NewAlertPeersCount(moniker, network, diff string) *Alert {
	i := Importance.Info
	m := fmt.Sprintf("Peers count of '%s' , \nNetwork: %s \n%s", moniker, network, diff)
	return New(i, m)
}

func NewAlertCatchingUp(moniker, network string) *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Node '%s'\n Net:%s\n Catching up", moniker, network)
	return New(i, m)
}

func NewAlertBlocksDelta(moniker string, diff int64) *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Node '%s' %v blocks behind", moniker, diff)
	return New(i, m)
}

func NewAlertSynced(moniker string) *Alert {
	i := Importance.OK
	m := fmt.Sprintf("Node '%s' Synced", moniker)
	return New(i, m)
}

func NewAlertAccessDelays(moniker, network, rpc string) *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Experiencing delays when trying to access '%s' node. \nNet: %s , \nNode RPC: %s",
		moniker,
		network,
		rpc,
	)
	return New(i, m)
}

func NewAlertAccessRestored(moniker,
	network string,
	lastSeenAt time.Time,
	timeDelta time.Duration) *Alert {
	i := Importance.OK
	m := fmt.Sprintf(
		"Node '%s'\nNet: %s\n is now accessible.\nNode became inaccessible at %s and was inaccessible for (at most) %s",
		moniker,
		network,
		lastSeenAt,
		timeDelta,
	)
	return New(i, m)
}

func NewAlertGithubRelease(tagName, repoName, releaseDesc string) *Alert {
	i := Importance.GH
	m := fmt.Sprintf("Release %s of %s has just been released.\n%s",
		tagName,
		repoName,
		releaseDesc)
	return New(i, m)
}
