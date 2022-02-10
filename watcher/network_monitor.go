package watcher

import (
	"encoding/base64"
	"log"
	"sync"
	"time"

	"github.com/gadost/telescope/conf"
	"github.com/gadost/telescope/event"
)

var wgNetMonitor sync.WaitGroup

type Participation struct {
	CheckedBlock               int64
	CountMissedSignatureInARow int
	Alert                      bool
}

type ValidatorInfo struct {
	Validators []Validator
}
type Validator struct {
	Address       string
	Moniker       string
	Network       string
	Participation Participation
}

func Proposal() {

}

// BlockProducingParticipation check blocks for our sign
func BlockProducingParticipation(cfg conf.ChainsConfig, chains []string) {
	for _, chainName := range chains {
		*node = cfg.Chain[chainName]

		for n, nodeConf := range node.Node {
			if nodeConf.NetworkMonitoringEnabled {
				log.Printf("Starting network monitoring for %s at %s ", chainName, nodeConf.RPC)
				validatorInfo := FindValidators(n, chainName)
				if len(validatorInfo.Validators) != 0 {
					wgNetMonitor.Add(1)
					go Scan(n, chainName, validatorInfo)
				}
			}
		}
	}
	wgNetMonitor.Wait()
}

// FindValidators find validators in configs
func FindValidators(n int, chainName string) *ValidatorInfo {
	// Need some sleep for preventing panic during bootstrap
	time.Sleep(10 * time.Second)
	var validatorInfo = new(ValidatorInfo)
	valPos := 0
	for _, f := range Chains.Chain[chainName].Node {
		if f.Role == "validator" {
			validatorInfo.Validators = append(validatorInfo.Validators, Validator{})
			validatorInfo.Validators[valPos].Moniker = f.Status.NodeInfo.Moniker
			validatorInfo.Validators[valPos].Network = f.Status.NodeInfo.Network
			validatorInfo.Validators[valPos].Address = base64.StdEncoding.EncodeToString(f.Status.ValidatorInfo.Address)
			valPos += 1
		}
	}
	return validatorInfo
}

// Scan for block participating
func Scan(n int, chainName string, v *ValidatorInfo) {
	client := Chains.Chain[chainName].Node[n].Client
	for {
		health, err := client.Health(ctx)
		if err != nil {
			log.Printf("Health: %s : Experiencing connection troubles %+v", Chains.Chain[chainName].Node[n].RPC, health)
		} else {
			if v.Validators[0].Participation.CheckedBlock == 0 {
				status, err := client.Status(ctx)
				if err == nil {
					v.Validators[0].Participation.CheckedBlock = status.SyncInfo.LatestBlockHeight
				}

			}
			commit, _ := client.Commit(ctx, &v.Validators[0].Participation.CheckedBlock)
			if commit != nil {
				for nA, vA := range v.Validators {

					s := len(commit.Commit.Signatures)
					t := len(commit.Commit.Signatures)
					for _, k := range commit.Commit.Signatures {
						if vA.Address == base64.StdEncoding.EncodeToString(k.ValidatorAddress) {
							t -= 1
						}
					}
					if s == t {
						v.Validators[nA].Participation.CountMissedSignatureInARow += 1
					} else {
						v.Validators[nA].Participation.CountMissedSignatureInARow = 0
						v.Validators[nA].Participation.Alert = false
					}
					if v.Validators[nA].Participation.CountMissedSignatureInARow ==
						int(Chains.Chain[chainName].Info.BlocksMissedInARow) &&
						!v.Validators[nA].Participation.Alert &&
						Chains.Chain[chainName].Info.BlocksMissedInARow != 0 {
						event.BlockMissedTracker(v.Validators[nA].Moniker,
							v.Validators[nA].Network,
							v.Validators[nA].Participation.CountMissedSignatureInARow,
						)
						v.Validators[nA].Participation.Alert = true
					}
				}
			}

		}
		time.Sleep(10 * time.Second)
	}

}
