package conf

import (
    "log"
    "os"

    "github.com/gadost/telescope/errors"

    "github.com/BurntSushi/toml"
)

type Config struct {
    Chains map[string]chain
}

type chain struct {
    Name  string
    Nodes map[string]node
}

type node struct {
    Role                  string
    Address               string
    NetworkMonitorEnabled bool
}

var conf Config

func ConfExist() {
    f := "~/.telescope/config.toml"
    if _, err := os.Stat(f); err != nil {
        log.Printf(errors.ConfNotFound, err)
    }

}

func ConfLoad(data string) {
    _, err := toml.Decode(data, &conf)
    if err != nil {
        panic(errors.ConfInvalid)
    }
}
