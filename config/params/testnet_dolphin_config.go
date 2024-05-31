package params

// UseDolphinNetworkConfig uses the Dolphin beacon chain specific network config.
func UseDolphinNetworkConfig() {
	cfg := BeaconNetworkConfig().Copy()
	cfg.BootstrapNodes = []string{
		// Dolphin testnet boot nodes
		"enr:-LG4QP-_uY01P5mD4-RrJzpuZlvLpyedCtkMU0q5Le69m60OY6sE5Bz0VI_592ujVjvxuvaJe_twOrpbq4GW0hTpXVGGAY_NUtwah2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEIkBUf4RvdmVykPWl_UIAAAAA__________-Jc2VjcDI1NmsxoQOp3abawRE7r3d9Kjn7qKrBzMMqc5p8AVpb7a81_MYkcYN1ZHCCyyA",
	}
	OverrideBeaconNetworkConfig(cfg)
}

// DolphinConfig defines the config for the Dolphin beacon chain testnet.
func DolphinConfig() *BeaconChainConfig {
	cfg := MainnetConfig().Copy()
	cfg.ConfigName = DolphinName
	cfg.PresetBase = "dolphin"
	cfg.DepositChainID = 541761
	cfg.DepositNetworkID = 541761
	cfg.GenesisForkVersion = []byte{0x00, 0x00, 0x00, 0x28}
	cfg.AltairForkEpoch = 2
	cfg.AltairForkVersion = []byte{0x01, 0x00, 0x00, 0x28}
	cfg.BellatrixForkEpoch = 4
	cfg.BellatrixForkVersion = []byte{0x02, 0x00, 0x00, 0x28}
	cfg.CapellaForkEpoch = 10
	cfg.CapellaForkVersion = []byte{0x03, 0x00, 0x00, 0x28}
	cfg.InitializeForkSchedule()
	return cfg
}
