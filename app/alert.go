package app

import (
	"fmt"
	"sync"
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
var alertSystems = &MainConfig
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

func (e *Event) NewAlertBlockMissed() *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Validator '%s' , \nNetwork: %s \nSignature missed  in last %v blocks in a row",
		e.Moniker,
		e.Network,
		e.Missed,
	)
	return New(i, m)
}

func (e *Event) NewAlertVotingPower() *Alert {
	i := Importance.Info
	m := fmt.Sprintf("Voting Power of '%s' , \nNetwork: %s \n%s", e.Moniker, e.Network, e.Diff)
	return New(i, m)
}

func (e *Event) NewAlertPeersCount() *Alert {
	i := Importance.Info
	m := fmt.Sprintf("Peers count of '%s' , \nNetwork: %s \n%s", e.Moniker, e.Network, e.Diff)
	return New(i, m)
}

func (e *Event) NewAlertCatchingUp() *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Node '%s'\n Net:%s\n Catching up", e.Moniker, e.Network)
	return New(i, m)
}

func (e *Event) NewAlertBlocksDelta() *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Node '%s' %v blocks behind", e.Moniker, e.Diff)
	return New(i, m)
}

func (e *Event) NewAlertSynced() *Alert {
	i := Importance.OK
	m := fmt.Sprintf("Node '%s' Synced", e.Moniker)
	return New(i, m)
}

func (e *Event) NewAlertAccessDelays() *Alert {
	i := Importance.Urgent
	m := fmt.Sprintf("Experiencing delays when trying to access '%s' node. \nNet: %s , \nNode RPC: %s",
		e.Moniker,
		e.Network,
		e.RPC,
	)
	return New(i, m)
}

func (e *Event) NewAlertAccessRestored() *Alert {
	i := Importance.OK
	m := fmt.Sprintf(
		"Node '%s'\nNet: %s\n is now accessible.\nNode became inaccessible at %s and was inaccessible for (at most) %s",
		e.Moniker,
		e.Network,
		e.LastSeenAt,
		e.TimeDelta,
	)
	return New(i, m)
}

func (e *Event) NewAlertGithubRelease() *Alert {
	i := Importance.GH
	m := fmt.Sprintf("Release %s of %s has just been released.\n%s",
		e.TagName,
		e.RepoName,
		e.ReleaseDesc)
	return New(i, m)
}
