package conf

import (
    "log"
    "os"
    "time"

    "github.com/gadost/telescope/errors"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Title   string
    Owner   ownerInfo
    DB      database `toml:"database"`
    Servers map[string]server
    Clients clients
}

type ownerInfo struct {
    Name string
    Org  string `toml:"organization"`
    Bio  string
    DOB  time.Time
}

type database struct {
    Server  string
    Ports   []int
    ConnMax int `toml:"connection_max"`
    Enabled bool
}

type server struct {
    IP string
    DC string
}

type clients struct {
    Data  [][]interface{}
    Hosts []string
}

var conf Config

func ConfigExist() {
    f := "~/.telescope/example.toml"
    if _, err := os.Stat(f); err != nil {
        log.Printf(errors.ConfNotFound, err)
    }

}

func ConfigValid(data string) {
    _, err := toml.Decode(data, &conf)
    if err != nil {
        panic(errors.ConfInvalid)
    }
}
