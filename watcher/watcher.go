package watcher

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/event"
	tmint "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/rpc/coretypes"
	tele "gopkg.in/telebot.v3"
)

var ctx = context.TODO()
var node = &conf.Nodes{}
var wgWatcher sync.WaitGroup
var Node conf.Node
var Chains = &conf.ChainsConfig{}

// Split nodes and run gorutine per node
func ThreadsSplitter(cfg conf.ChainsConfig, chains []string) {
	Chains = &cfg

	for _, chainName := range chains {
		*node = cfg.Chain[chainName]

		for n, nodeConf := range node.Node {
			if nodeConf.MonitoringEnabled {
				log.Printf("Starting monitoring for %s", nodeConf.RPC)
				wgWatcher.Add(1)
				go Thread(nodeConf.RPC, chainName, n)
			}
		}
	}

	wgWatcher.Wait()
}

func Thread(rpc string, chainName string, n int) {
	defer wgWatcher.Done()

	client, err := tmint.New(rpc)
	if err != nil {
		log.Println(rpc, err)
	} else {
		Chains.Chain[chainName].Node[n].Client = client
		err = client.Start()
		if err != nil {
			log.Println(rpc, err)
		}

		var counter int

		for {
			health, _ := client.Health(ctx)
			if health == nil {
				counter += 1
				log.Printf("Health: %s : Experiencing connection troubles %+v", rpc, health)
			} else {
				counter = 0

				status, _ := client.Status(ctx)
				if status != nil {
					CheckStatus(status, chainName, rpc)
				}

				netInfo, _ := client.NetInfo(ctx)
				if netInfo != nil {
					CheckPeers(netInfo, chainName, rpc)
				}
			}

			counter = CheckHealth(chainName, rpc, counter)

			time.Sleep(10 * time.Second)
		}
	}
}

func CheckPeers(res *coretypes.ResultNetInfo, chainName string, rpc string) {
	for i, n := range Chains.Chain[chainName].Node {
		if n.RPC == rpc {
			status := Chains.Chain[chainName].Node[i].Status
			if status.BootstrappedNetInfo {
				event.PeersCount(
					status.PeersCount,
					res.NPeers,
					status.NodeInfo.Moniker,
					status.NodeInfo.Network,
				)
			} else {
				Chains.Chain[chainName].Node[i].Status.BootstrappedNetInfo = true
			}

			Chains.Chain[chainName].Node[i].Status.PeersCount = res.NPeers
		}
	}
}

func CheckHealth(chainName string, rpc string, counter int) int {
	for i, n := range Chains.Chain[chainName].Node {
		if n.RPC == rpc {
			status := Chains.Chain[chainName].Node[i].Status

			if counter == 5 {
				Chains.Chain[chainName].Node[i].Status.HealthStateBad = true
				Chains.Chain[chainName].Node[i].Status.LastSeenAt = Chains.Chain[chainName].Node[i].Status.SyncInfo.LatestBlockTime
			}

			_, resolved := event.HealthCheck(
				status.NodeInfo.Moniker,
				status.NodeInfo.Network,
				rpc,
				counter,
				time.Now().UTC().Sub(status.LastSeenAt.UTC()),
				status.LastSeenAt,
				status.HealthStateBad,
			)

			if resolved {
				Chains.Chain[chainName].Node[i].Status.HealthStateBad = false
			}
		}
	}

	return counter
}

func CheckStatus(res *coretypes.ResultStatus, chainName string, rpc string) {

	for i, n := range Chains.Chain[chainName].Node {
		if n.RPC == rpc {
			status := Chains.Chain[chainName].Node[i].Status

			status.NodeInfo = res.NodeInfo
			if status.BootstrappedStatus {
				event.VotingPower(
					status.ValidatorInfo.VotingPower,
					res.ValidatorInfo.VotingPower,
					Chains.Chain[chainName].Info.VotingPowerChanges,
					status.NodeInfo.Moniker,
					status.NodeInfo.Network,
				)
			} else {
				Chains.Chain[chainName].Node[i].Status.BootstrappedStatus = true
			}

			event.CatchingUpState(
				Chains.Chain[chainName].Node[i].Status.SyncInfo.CatchingUp,
				res.SyncInfo.CatchingUp,
				status.NodeInfo.Moniker,
				status.NodeInfo.Network,
				SearchLatestBlock(chainName)-Chains.Chain[chainName].Node[i].Status.SyncInfo.LatestBlockHeight,
				Chains.Chain[chainName].Info.BlocksMissedInARow,
			)
			Chains.Chain[chainName].Node[i].Status.NodeInfo = res.NodeInfo
			Chains.Chain[chainName].Node[i].Status.SyncInfo = res.SyncInfo
			Chains.Chain[chainName].Node[i].Status.ValidatorInfo = res.ValidatorInfo
		}
	}
}

func SearchLatestBlock(chainName string) int64 {
	var lBs []int64
	for _, n := range Chains.Chain[chainName].Node {
		lBs = append(lBs, n.Status.SyncInfo.LatestBlockHeight)
	}
	var lB = lBs[0]
	if len(lBs) > 1 {

		for _, v := range lBs {
			if v > lB {
				lB = v
			}
		}
	}
	return lB
}
func StatusCollection() string {
	var collection string
	var cu string
	var badRPC []string
	collection = collection + "*Status:*\n\n"

	for i := range Chains.Chain {
		for _, k := range Chains.Chain[i].Node {
			if k.Status.SyncInfo.LatestBlockHeight > 0 {
				collection += "*Net:* `" + k.Status.NodeInfo.Network + "`\n*Moniker:* `" + k.Status.NodeInfo.Moniker + "`\n"
				if k.Status.SyncInfo.CatchingUp {
					cu = "Yes"
				} else {
					cu = "No"
				}

				collection += "*CatchingUp:* `" + cu + "`\n"
				collection += "*Last known height:* `" + strconv.Itoa(int(k.Status.SyncInfo.LatestBlockHeight)) + "`\n"
				collection += "*Last known block time :* `" + k.Status.SyncInfo.LatestBlockTime.Format("2006-01-02 15:04:05") + "`\n"
				collection += "`_________________________`\n"
			} else {
				if k.MonitoringEnabled {
					badRPC = append(badRPC, k.RPC)
				}
			}
		}
	}

	if len(badRPC) > 0 {
		collection += "*🔴Unreachable RPCs:*\n`" + strings.Join(badRPC, "\n") + "`"
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
		return c.Send(StatusCollection(), "MarkdownV2")
	})

	b.Start()
}
