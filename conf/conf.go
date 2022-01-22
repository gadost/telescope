package conf

import (
	"flag"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var configPath string
var userHome = os.Getenv("HOME")
var cfg Config
var infra nodes
var chains []string
var zNodes = &nodes{}

func (n *nodes) Reset() {
	*n = *zNodes
}
func init() {

	flag.StringVar(&configPath, "confd", userHome+"/.telescope/conf.d", "path to configs dir")
	flag.Parse()

}

//Chain struct for <chain>.toml configs
type Config struct {
	Chain map[string]nodes
}

type nodes struct {
	Info info
	Node []node
}
type info struct {
	Mainnet  bool
	Telegram bool
}

type node struct {
	Role                     string
	Address                  string
	NetworkMonitoringEnabled bool `toml:"network_monitoring_enabled"`
}

//Check existence of confd folder
func ConfLoad() (Config, []string) {
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(configPath)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		buildConf(files)
	}
	return cfg, chains
}

// buildConf load parsed config to struct
func buildConf(files []fs.FileInfo) {

	for _, f := range files {
		if _, err := toml.DecodeFile(configPath+"/"+f.Name(), &infra); err != nil {
			log.Fatal(err)
		}
		log.Println(infra)
		//prevent panic on nil map
		if cfg.Chain == nil {
			cfg.Chain = make(map[string]nodes)
		}
		cfg.Chain[fileNameWithoutExtSliceNotation(f.Name())] = infra
		chains = append(chains, fileNameWithoutExtSliceNotation(f.Name()))
		// Should be reseted, because using one univerasl struct
		infra.Reset()
	}
}

// Remove extesion  for add to chains slice
func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
