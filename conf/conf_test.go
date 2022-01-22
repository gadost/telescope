package conf

import (
	"fmt"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestDebugReturn(t *testing.T) {
	cfg, chains := ConfLoad()
	fmt.Println(cfg)
	fmt.Println(chains)
	fmt.Println(chains[0])
	fmt.Println(cfg.Chain["gravity"].Node[0].Role)
	fmt.Println(cfg.Chain["gravity"].Info.Mainnet)
}
