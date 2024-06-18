//go:build !noMainnetGenesis
// +build !noMainnetGenesis

package genesis

import (
	_ "embed"
)

var (
	//go:embed mainnet.ssz.snappy
	mainnetRawSSZCompressed []byte // 1.8Mb
)

func init() {
	//embeddedStates[params.MainnetName] = &mainnetRawSSZCompressed
}
