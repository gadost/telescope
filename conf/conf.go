package conf

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/coretypes"
	"github.com/tendermint/tendermint/types"
)

var (
	chains         []string
	MainConfig     Config
	chainsConfig   ChainsConfig
	ConfdPath      string
	infra          Nodes
	mainConfigName = "telescope.toml"
	userHome       = os.Getenv("HOME")
	zNodes         = &Nodes{}
)

//Chain struct for <chain>.toml configs
type ChainsConfig struct {
	Chain map[string]Nodes
}

type Nodes struct {
	Info struct {
		Github             string `toml:"github"`
		Mainnet            bool
		Telegram           bool
		VotingPowerChanges int64 `toml:"voting_power_changes"`
		BlocksMissedInARow int64 `toml:"blocks_missed_in_a_row"`
		PeersCount         int64 `toml:"peers_count"`
	}
	Node []Node
}
type Node struct {
	Role                     string
	RPC                      string `toml:"rpc"`
	NetworkMonitoringEnabled bool   `toml:"network_monitoring_enabled"`
	MonitoringEnabled        bool   `toml:"monitoring_enabled"`
	Status                   NodeStatus
}

type NodeStatus struct {
	NodeInfo            types.NodeInfo
	SyncInfo            coretypes.SyncInfo
	HealthStateBad      bool
	BootstrappedStatus  bool
	BootstrappedNetInfo bool

	ValidatorInfo coretypes.ValidatorInfo

	BlockMissedTracker uint64
	PeersCount         int

	LastSeenProblemsAt time.Time
	LastSeenAt         time.Time
}

type NodeSyncInfo struct {
	LatestBlockHash   bytes.HexBytes
	LatestBlockHeight int64
	CatchingUp        bool
	LatestBlockTime   time.Time
}

// telescope.toml
type Config struct {
	Settings struct {
		DowntimeInterval     int  `toml:"downtime_interval"`
		GithubReleaseMonitor bool `toml:"github_release_monitor"`
	} `toml:"settings"`
	Telegram struct {
		Enabled bool   `toml:"enabled"`
		Token   string `toml:"token"`
		ChatID  string `toml:"chat_id"`
	} `toml:"telegram"`
	Discord struct {
		Enabled bool `toml:"enabled"`
	} `toml:"discord"`
	Twilio struct {
		Enabled bool `toml:"enabled"`
	} `toml:"twilio"`
	Mail struct {
		Enabled bool `toml:"enabled"`
	} `toml:"mail"`
	Sms struct {
		Enabled bool `toml:"enabled"`
	} `toml:"sms"`
}

func (n *Nodes) Reset() {
	*n = *zNodes
}

func init() {
	flag.StringVar(&ConfdPath, "confd", userHome+"/.telescope/conf.d", "path to configs dir")
	flag.Parse()
}

//Check existence of confd folder
func ConfLoad() (ChainsConfig, []string) {
	if _, err := os.Stat(ConfdPath); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(ConfdPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		buildConf(files)
	}
	return chainsConfig, chains
}

// buildConf load parsed config to struct
func buildConf(files []fs.FileInfo) {

	for _, f := range files {
		switch f.Name() {
		case mainConfigName:
			if _, err := toml.DecodeFile(ConfdPath+"/"+f.Name(), &MainConfig); err != nil {
				log.Fatal(err)
			}
		default:

			if _, err := toml.DecodeFile(ConfdPath+"/"+f.Name(), &infra); err != nil {
				log.Fatal(err)
			}
			//prevent panic on nil map
			if chainsConfig.Chain == nil {
				chainsConfig.Chain = make(map[string]Nodes)
			}
			chainsConfig.Chain[fileNameWithoutExt(f.Name())] = infra
			chains = append(chains, fileNameWithoutExt(f.Name()))
			// Should be reseted, because using one univerasl struct
			infra.Reset()
		}
	}
}

// Remove extesion  for add to chains slice
func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
