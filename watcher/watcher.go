package watcher

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gadost/telescope/conf"
	tmint "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/coretypes"
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
					// ALERT HERE
					fmt.Println(s.ValidatorInfo.VotingPower)
				} else if s.ValidatorInfo.VotingPower-Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower >= Chains.Chain[c].Info.VotingPowerChanges {
					// ALERT HERE
					fmt.Println(s.ValidatorInfo.VotingPower)
				}
			}
			Chains.Chain[c].Node[i].Status.ValidatorInfo.VotingPower = s.ValidatorInfo.VotingPower

			// CATCHINGUP STATE
			switch Chains.Chain[c].Node[i].Status.SyncInfo.CatchingUp {
			case false:
				switch s.SyncInfo.CatchingUp {
				case true:
					Chains.Chain[c].Node[i].Status.SyncInfo.CatchingUp = s.SyncInfo.CatchingUp
					//ALERT HERE CATCHING
					if Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockHeight-s.SyncInfo.LatestBlockHeight > Chains.Chain[c].Info.BlocksMissedInARow {
						// ALERT BLOCKS TO GO
						fmt.Println(Chains.Chain[c].Node[i].Status.SyncInfo.LatestBlockHeight - s.SyncInfo.LatestBlockHeight)
					}
				}
			case true:
				switch s.SyncInfo.CatchingUp {
				case false:
					Chains.Chain[c].Node[i].Status.SyncInfo.CatchingUp = s.SyncInfo.CatchingUp
					// ALERT HERE DONE
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
					//ALERT HERE DECREASED
					fmt.Println(ni.NPeers)
				} else if Chains.Chain[c].Node[i].Status.PeersCount-ni.NPeers < 0 {
					//ALERT HERE INCREASED
					fmt.Println(ni.NPeers)
				}
			}
			Chains.Chain[c].Node[i].Status.PeersCount = ni.NPeers
		}
	}
}
