package event

import (
	"fmt"
	"time"

	"github.com/gadost/telescope/alert"
)

type Context struct {
	VotingPower bool
	PeersCount  bool
}

// BlockMissedTracker check missed blocks in a row
func BlockMissedTracker(moniker, network string, missed int) {
	alert.New(alert.Importance.Urgent,
		fmt.Sprintf("Validator '%s' , \nNetwork: %s \n Signature missed  in last %v blocks in a row", moniker, network, missed),
	)
}

// Difference diff two values , comparing with minimal difference and return increase/decrease message for alert
func Difference(one, two, min int64, ctx Context) (string, bool) {
	switch d := one - two; {
	case d > 0 && ((ctx.VotingPower && d > min) || (ctx.PeersCount && two < min)):
		return fmt.Sprintf("DECREASED by %v \nfrom %v to %v", d, one, two), true
	case d < 0 && ((ctx.VotingPower && d*(-1) > min) || (ctx.PeersCount && one < min && two < 10)):
		return fmt.Sprintf("INCREASED by %v \nfrom %v to %v", d*(-1), one, two), true
	default:
		return "", false
	}
}

// VotingPower comparing voting power of validator , if changed - send alert
func VotingPower(sVP, rVP, sVPC int64, moniker, network string) {
	diff, changed := Difference(sVP, rVP, sVPC, Context{VotingPower: true})
	if changed {
		alert.New(alert.Importance.Info,
			fmt.Sprintf("Voting Power of '%s' , \nNetwork: %s \n%s", moniker, network, diff),
		)
	}
}

// PeersCount check for peers count, alert if changed
func PeersCount(sPC, rPC int, moniker, network string) {
	diff, changed := Difference(int64(sPC), int64(rPC), 10, Context{PeersCount: true})
	if changed {
		alert.New(alert.Importance.Info,
			fmt.Sprintf("Peers count of '%s' , \nNetwork: %s \n%s", moniker, network, diff),
		)
	}
}

// CatchingUpState check cachingUp state , alert if true
func CatchingUpState(sCU, rCU bool, moniker, network string, diff, maxDiff int64) {
	switch sCU {
	case false:
		switch rCU {
		case true:
			alert.New(
				alert.Importance.Urgent,
				fmt.Sprintf("Node '%s'\n Net:%s\n Catching up", moniker, network),
			)
			if diff > maxDiff {
				alert.New(
					alert.Importance.Urgent,
					fmt.Sprintf("Node '%s' %v blocks behind", moniker, diff),
				)
			}
		}
	case true:
		switch rCU {
		case false:
			alert.New(alert.Importance.OK,
				fmt.Sprintf("Node '%s' Synced", moniker),
			)
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
func HealthCheck(moniker, network, rpc string, counter int, timeDelta time.Duration, lastSeenAt time.Time, lastStatus bool) (bool, bool) {
	var resolved = false
	if counter == 5 {
		alert.New(
			alert.Importance.Urgent,
			fmt.Sprintf("Experiencing delays when trying to access '%s' node. \nNet: %s , \nNode RPC: %s",
				Unknown(moniker),
				Unknown(network),
				rpc,
			),
		)

		return true, resolved

	} else if counter > 5 {

		return true, resolved

	} else if counter == 0 && lastStatus {
		alert.New(
			alert.Importance.OK,
			fmt.Sprintf("Node '%s'\nNet: %s\n is now accessible.\nNode became inaccessible at %s and was inaccessible for (at most) %s",
				Unknown(moniker),
				Unknown(network),
				lastSeenAt,
				timeDelta,
			),
		)

		resolved = true
		return false, resolved

	} else if counter == 0 && !lastStatus {

		return false, resolved

	} else {

		return true, resolved

	}
}
