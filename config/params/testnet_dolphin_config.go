package params

// UseDolphinNetworkConfig uses the Dolphin beacon chain specific network config.
func UseDolphinNetworkConfig() {
	cfg := BeaconNetworkConfig().Copy()
	cfg.BootstrapNodes = []string{
		// Dolphin testnet boot nodes
		"enr:-Iq4QMCTfIMXnow27baRUb35Q8iiFHSIDBJh6hQM5Axohhf4b6Kr_cOCu0htQ5WvVqKvFgY28893DHAg8gnBAXsAVqmGAX53x8JggmlkgnY0gmlwhLKAlv6Jc2VjcDI1NmsxoQK6S-Cii_KmfFdUJL2TANL3ksaKUnNXvTCv1tLwXs0QgIN1ZHCCIyk",
		"enr:-KG4QE5OIg5ThTjkzrlVF32WT_-XT14WeJtIz2zoTqLLjQhYAmJlnk4ItSoH41_2x0RX0wTFIe5GgjRzU2u7Q1fN4vADhGV0aDKQqP7o7pAAAHAyAAAAAAAAAIJpZIJ2NIJpcISlFsStiXNlY3AyNTZrMaEC-Rrd_bBZwhKpXzFCrStKp1q_HmGOewxY3KwM8ofAj_ODdGNwgiMog3VkcIIjKA",
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
	cfg.AltairForkEpoch = 2
	cfg.AltairForkVersion = []byte{0x1, 0x00, 0x00, 0x10}
	cfg.BellatrixForkEpoch = 4
	cfg.BellatrixForkVersion = []byte{0x2, 0x00, 0x00, 0x10}
	cfg.CapellaForkEpoch = 10
	cfg.CapellaForkVersion = []byte{0x3, 0x00, 0x00, 0x10}
	cfg.InitializeForkSchedule()
	return cfg
}
