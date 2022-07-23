package app

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	ChainNamePrefix = "chain_name"
)

type LevelDB struct {
	*leveldb.DB
}

func NewDB(path string) *LevelDB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Panic(err)
	}
	return &LevelDB{db}
}

func (db *LevelDB) OnStartSetChains(chains []string) {
	for _, c := range chains {
		err := db.Put([]byte(ChainNamePrefix+c), nil, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func (db *LevelDB) SetGithubLatestRelease(owner string, repo string, version string) error {
	err := db.Put([]byte(owner+repo), []byte(version), nil)
	return err
}

func (db *LevelDB) GetGithubLatestRelease(owner string, repo string) ([]byte, error) {
	version, err := db.Get([]byte(owner+repo), nil)
	return version, err
}
