package watcher

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gadost/telescope/alert"
	"github.com/gadost/telescope/conf"
	tmint "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/coretypes"
	tele "gopkg.in/telebot.v3"
)

var ctx = context.TODO()
var node = &conf.Nodes{}
var wg sync.WaitGroup
var Node conf.Node

//var Chains = &conf.ChainsConfig{}
var Chains = &conf.ChainsConfig{}

func ThreadsSplitter(cfg conf.ChainsConfig, chains []string) {
	Chains = &cfg
	for _, c := range chains {
		*node = cfg.Chain[c]
		for _, n := range node.Node {
			if n.NetworkMonitoringEnabled {
				wg.Add(1)
				go Thread(n.RPC, c)
			}
		}
	}
	wg.Wait()
}

func Thread(r string, c string) {
	defer wg.Done()
	t, err := tmint.New(r)

	if err != nil {
		log.Println(err)
	} else {
		err = t.Start(ctx)
		if err != nil {
			log.Println("Will trying reconnect later", err)
		}
		for {
			s, _ := t.Status(ctx)
			if s != nil {
				ParseStatus(s, c, r)
			}
			ni, _ := t.NetInfo(ctx)
			if ni != nil {
				ParseNetInfo(ni, c, r)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func ParseStatus(s *coretypes.ResultStatus, c string, r string) {
	for i, n := range Chains.Chain[c].Node {
		if n.RPC == r {
			log.Printf("%+v", Chains.Chain[c].Node[i].Status)
			log.Println(s.NodeInfo.Moniker)

			//Probably change to struct load ( just now mismatch types)
			Chains.Chain[c].Node[i].Status.NodeInfo.Moniker = s.NodeInfo.Moniker
			Chains.Chain[c].Node[i].Status.NodeInfo.Network = s.NodeInfo.Network
			Chains.Chain[c].Node[i].Status.NodeInfo.NodeID = string(s.NodeInfo.NodeID)
			Chains.Chain[c].Node[i].Status.ValidatorInfo.PubKey = s.ValidatorInfo.PubKey
			//VALIDATOR POWER CHANGES
			if Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower != s.ValidatorInfo.VotingPower {
				if Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower-s.ValidatorInfo.VotingPower >= Chains.Chain[c].Info.VotingPowerChanges {
					alert.New(alert.Importance.Info, fmt.Sprintf("Voting Power of '%s' , Network: %s DECREASED by %v to %v",
						Chains.Chain[c].Node[i].Status.NodeInfo.Moniker,
						Chains.Chain[c].Node[i].Status.NodeInfo.Network,
						Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower-s.ValidatorInfo.VotingPower,
						s.ValidatorInfo.VotingPower))
				} else if s.ValidatorInfo.VotingPower-Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower >= Chains.Chain[c].Info.VotingPowerChanges {
					alert.New(alert.Importance.Info, fmt.Sprintf("Voting Power of '%s' , Network: %s INCREASED by %v to %v",
						Chains.Chain[c].Node[i].Status.NodeInfo.Moniker,
						Chains.Chain[c].Node[i].Status.NodeInfo.Network,
						s.ValidatorInfo.VotingPower-Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower,
						s.ValidatorInfo.VotingPower))
				}
			}
			Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower = s.ValidatorInfo.VotingPower

			// CATCHINGUP STATE
			switch Chains.Chain[c].Node[i].Status.SyncInfo.CatchingUp {
			case false:
				switch s.SyncInfo.CatchingUp {
				case true:
					Chains.Chain[c].Node[i].Status.SyncInfo.CatchingUp = s.SyncInfo.CatchingUp
					alert.New(alert.Importance.Urgent, fmt.Sprintf("Node '%s' CatchingUp", Chains.Chain[c].Node[i].Status.NodeInfo.Moniker))
					if Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockHeight-s.SyncInfo.LatestBlockHeight > Chains.Chain[c].Info.BlocksMissedInARow {
						alert.New(alert.Importance.Urgent, fmt.Sprintf("Node '%s' %v blocks behind",
							Chains.Chain[c].Node[i].Status.NodeInfo.Moniker,
							Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockHeight-s.SyncInfo.LatestBlockHeight))
					}
				}
			case true:
				switch s.SyncInfo.CatchingUp {
				case false:
					Chains.Chain[c].Node[i].Status.SyncInfo.CatchingUp = s.SyncInfo.CatchingUp
					alert.New(alert.Importance.OK, fmt.Sprintf("Node '%s' Synced", Chains.Chain[c].Node[i].Status.NodeInfo.Moniker))
				}
			}

			Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockHash = s.SyncInfo.LatestBlockHash
			Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockHeight = s.SyncInfo.LatestBlockHeight
			Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockTime = s.SyncInfo.LatestBlockTime
		}
	}
}

func ParseNetInfo(ni *coretypes.ResultNetInfo, c string, r string) {
	for i, n := range Chains.Chain[c].Node {
		if n.RPC == r {
			if ni.NPeers <= 10 {
				if Chains.Chain[c].Node[i].Status.PeersCount-ni.NPeers > 0 {
					alert.New(alert.Importance.Urgent, fmt.Sprintf("Peers count DECREASED by %v to %v", Chains.Chain[c].Node[i].Status.PeersCount-ni.NPeers, ni.NPeers))
				} else if Chains.Chain[c].Node[i].Status.PeersCount-ni.NPeers < 0 {
					alert.New(alert.Importance.OK, fmt.Sprintf("Peers count INCREASED by %v to %v", ni.NPeers-Chains.Chain[c].Node[i].Status.PeersCount, ni.NPeers))
				}
			}
			Chains.Chain[c].Node[i].Status.PeersCount = ni.NPeers
		}
	}
}

func StatusCollection() string {
	var collection string
	var cu string
	collection = collection + "<b>Status:</b>\n\n"
	for i := range Chains.Chain {
		for _, k := range Chains.Chain[i].Node {
			collection += "<b>Net:</b> " + k.Status.NodeInfo.Network + "\n<b>Moniker:</b> " + k.Status.NodeInfo.Moniker + "\n"
			if k.Status.SyncInfo.CatchingUp {
				cu = "Yes"
			} else {
				cu = "No"
			}

			collection += "<b>CatchingUp:</b> " + cu + "\n"
			collection += "<b>Last known height:</b> " + strconv.Itoa(int(k.Status.SyncInfo.LatestBlockHeight)) + "\n"
			collection += "<b>Last seen at:</b> " + k.Status.SyncInfo.LatestBlockTime.Format("2006-01-02 15:04:05") + "\n"
			collection += "_________________________\n"
		}
	}
	return collection
}

func TelegramHandler() {
	var pref = tele.Settings{
		Token:  conf.MainConfig.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/status", func(c tele.Context) error {
		return c.Send(StatusCollection(), "HTML")
	})

	b.Start()
}
