package app

import (
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestConfig(t *testing.T) {
	cfg, chains := ConfLoad(UserHome)
	t.Log(chains)
	t.Log(cfg)

}
