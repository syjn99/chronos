package params

// UseDolphinNetworkConfig uses the Dolphin beacon chain specific network config.
func UseDolphinNetworkConfig() {
	cfg := BeaconNetworkConfig().Copy()
	cfg.BootstrapNodes = []string{
		// Dolphin testnet boot nodess
		"enr:-LG4QBg_EHDI1NooVMFiODW81B6h5Xw1xGBQT8VfsR_IBZgdFZaGrmyFUomH-kd9gWOqaf_tlseyLRSXMZBsy2vMPWmGAZBJSMwth2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEIkBUf4RvdmVykNBNsU8AAAAY__________-Jc2VjcDI1NmsxoQOaEHgEJPUhVaDcltNoWYFRCE7GcZccocBv6EW6uX1M2IN1ZHCCyyA",
	}
	OverrideBeaconNetworkConfig(cfg)
}

// DolphinConfig defines the config for the Dolphin beacon chain testnet.
func DolphinConfig() *BeaconChainConfig {
	cfg := MainnetConfig().Copy()
	cfg.ConfigName = DolphinName
	cfg.PresetBase = "dolphin"
	cfg.DepositChainID = 541762
	cfg.DepositNetworkID = 541762
	cfg.GenesisForkVersion = []byte{0x00, 0x00, 0x00, 0x28}
	cfg.AltairForkEpoch = 0
	cfg.AltairForkVersion = []byte{0x01, 0x00, 0x00, 0x28}
	cfg.BellatrixForkEpoch = 0
	cfg.BellatrixForkVersion = []byte{0x02, 0x00, 0x00, 0x28}
	cfg.CapellaForkEpoch = 4
	cfg.CapellaForkVersion = []byte{0x03, 0x00, 0x00, 0x28}
	cfg.IssuanceRate = [11]uint64{20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 0}
	cfg.MaxBoostYield = [11]uint64{0, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000, 10000000000}
	cfg.InitializeForkSchedule()
	return cfg
}
