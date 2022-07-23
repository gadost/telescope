package app

import (
	"fmt"
	"time"
)

type Event struct {
	Moniker     string
	Network     string
	Missed      int
	Diff        interface{}
	MaxDiff     int64
	RPC         string
	LastSeenAt  time.Time
	TimeDelta   time.Duration
	TagName     string
	RepoName    string
	ReleaseDesc string
}

func VotingPowerChanges(one, two, min int64) (string, bool) {
	switch d := one - two; {
	case d > 0 && d > min:
		return fmt.Sprintf("DECREASED by %v \nfrom %v to %v", d, one, two), true
	case d < 0 && d*(-1) > min:
		return fmt.Sprintf("INCREASED by %v \nfrom %v to %v", d*(-1), one, two), true
	default:
		return "", false
	}
}

func PeersCountChanges(one, two, min int64) (string, bool) {
	switch d := one - two; {
	case d > 0 && two < min:
		return fmt.Sprintf("DECREASED by %v \nfrom %v to %v", d, one, two), true
	case d < 0 && one < min && two < 10:
		return fmt.Sprintf("INCREASED by %v \nfrom %v to %v", d*(-1), one, two), true
	default:
		return "", false
	}
}

// VotingPower comparing voting power of validator , if changed - send alert
func VotingPower(sVP, rVP, sVPC int64, moniker, network string) {
	diff, changed := VotingPowerChanges(sVP, rVP, sVPC)
	if changed {
		e := Event{
			Moniker: moniker,
			Network: network,
			Diff:    diff,
		}
		e.NewAlertVotingPower().Send()
	}
}

// PeersCount check for peers count, alert if changed
func PeersCount(sPC, rPC int, moniker, network string) {
	diff, changed := PeersCountChanges(int64(sPC), int64(rPC), 10)
	if changed {
		e := Event{
			Moniker: moniker,
			Network: network,
			Diff:    diff,
		}
		e.NewAlertPeersCount().Send()
	}
}

// CatchingUpState check cachingUp state , alert if true
func CatchingUpState(sCU, rCU bool, moniker, network string, diff, maxDiff int64) {
	e := Event{
		Moniker: moniker,
		Network: network,
		Diff:    diff,
		MaxDiff: maxDiff,
	}
	switch sCU {
	case false:
		switch rCU {
		case true:
			e.NewAlertCatchingUp().Send()
			if diff > maxDiff {
				e.NewAlertBlocksDelta().Send()
			}
		}
	case true:
		switch rCU {
		case false:
			e.NewAlertSynced().Send()
		}
	}
}

// replace nil string with "UNKNOWN"
func Unknown(s string) string {
	if s == "" {
		return "UNKNOWN"
	} else {
		return s
	}
}

// HealthCheck check node health
func HealthCheck(moniker, network, rpc string, counter int,
	timeDelta time.Duration, lastSeenAt time.Time, lastStatus bool) (bool, bool) {
	e := Event{
		Moniker:    moniker,
		Network:    network,
		RPC:        rpc,
		LastSeenAt: lastSeenAt,
		TimeDelta:  timeDelta,
	}
	var resolved = false
	if counter == 5 {
		e.NewAlertAccessDelays().Send()
		return true, resolved
	} else if counter > 5 {
		return true, resolved
	} else if counter == 0 && lastStatus {
		e.NewAlertAccessRestored().Send()
		resolved = true
		return false, resolved
	} else if counter == 0 && !lastStatus {
		return false, resolved
	} else {
		return true, resolved
	}
}
