package params

import (
	"math"
)

// UseDolphinNetworkConfig uses the Dolphin beacon chain specific network config.
func UseDolphinNetworkConfig() {
	cfg := BeaconNetworkConfig().Copy()
	cfg.BootstrapNodes = []string{
		// Dolphin testnet boot nodes
		"enr:-LG4QBg_EHDI1NooVMFiODW81B6h5Xw1xGBQT8VfsR_IBZgdFZaGrmyFUomH-kd9gWOqaf_tlseyLRSXMZBsy2vMPWmGAZBJSMwth2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEIkBUf4RvdmVykNBNsU8AAAAY__________-Jc2VjcDI1NmsxoQOaEHgEJPUhVaDcltNoWYFRCE7GcZccocBv6EW6uX1M2IN1ZHCCyyA",
	}
	OverrideBeaconNetworkConfig(cfg)
}

// DolphinConfig defines the config for the Dolphin beacon chain testnet.
func DolphinConfig() *BeaconChainConfig {
	cfg := MainnetConfig().Copy()
	cfg.ConfigName = DolphinName
	cfg.GenesisForkVersion = []byte{0x0, 0x00, 0x00, 0x10}
	cfg.DepositChainID = 541761
	cfg.DepositNetworkID = 541761
	cfg.AltairForkEpoch = 0
	cfg.AltairForkVersion = []byte{0x1, 0x00, 0x00, 0x10}
	cfg.BellatrixForkEpoch = 0
	cfg.BellatrixForkVersion = []byte{0x2, 0x00, 0x00, 0x10}
	cfg.CapellaForkEpoch = 4
	cfg.CapellaForkVersion = []byte{0x3, 0x00, 0x00, 0x10}
	cfg.DenebForkEpoch = math.MaxUint64
	cfg.DenebForkVersion = []byte{0x4, 0x00, 0x00, 0x10}
	cfg.InitializeForkSchedule()
	return cfg
}
