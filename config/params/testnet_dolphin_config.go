package params

import (
	"math"
)

// UseDolphinNetworkConfig uses the Dolphin beacon chain specific network config.
func UseDolphinNetworkConfig() {
	cfg := BeaconNetworkConfig().Copy()
	cfg.BootstrapNodes = []string{
		// Dolphin testnet boot nodes
		"enr:-LG4QGaXDTDc5_-AvUXuWxoYlT2Ce9dSlLi4Kx0Wzv7PFBSFWqRubay-w-IY5lay30YpEbP6_yNQtXa1QcrRD1PSdYqGAZFLTRaKh2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEgMdLF4RvdmVykNBNsU8AAAAY__________-Jc2VjcDI1NmsxoQOr1euFU8IZdyGo8jbIzJD0Z8VcRnt9xrIF-aOrRvQjPYN1ZHCCyyA",
	}
	OverrideBeaconNetworkConfig(cfg)
}

// DolphinConfig defines the config for the Dolphin beacon chain testnet.
func DolphinConfig() *BeaconChainConfig {
	cfg := MainnetConfig().Copy()
	cfg.ConfigName = DolphinName
	cfg.GenesisForkVersion = []byte{0x0, 0x00, 0x00, 0x28}
	cfg.DepositChainID = 541764
	cfg.DepositNetworkID = 541764
	cfg.AltairForkEpoch = 0
	cfg.AltairForkVersion = []byte{0x1, 0x00, 0x00, 0x28}
	cfg.BellatrixForkEpoch = 0
	cfg.BellatrixForkVersion = []byte{0x2, 0x00, 0x00, 0x28}
	cfg.CapellaForkEpoch = 4
	cfg.CapellaForkVersion = []byte{0x3, 0x00, 0x00, 0x28}
	cfg.DenebForkEpoch = math.MaxUint64
	cfg.DenebForkVersion = []byte{0x4, 0x00, 0x00, 0x28}
	cfg.IssuanceRate = [11]uint64{20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 0}
	cfg.MaxBoostYield = [11]uint64{0, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000}
	cfg.InitializeForkSchedule()
	return cfg
}
