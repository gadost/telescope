package conf

import (
    "log"
    "os"
    "time"

    "github.com/gadost/telescope/errors"

    "github.com/BurntSushi/toml"
)

type Config struct {

}

var conf Config

func ConfigExist() {
    f := "~/.telescope/config.toml"
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
