package params

import (
	"testing"
)

// SetupTestConfigCleanup preserves configurations allowing to modify them within tests without any
// restrictions, everything is restored after the test.
func SetupTestConfigCleanup(t testing.TB) {
	prevDefaultBeaconConfig := mainnetBeaconConfig.Copy()
	temp := configs.getActive().Copy()
	undo, err := SetActiveWithUndo(temp)
	if err != nil {
		t.Error(err)
	}
	prevNetworkCfg := networkConfig.Copy()
	t.Cleanup(func() {
		mainnetBeaconConfig = prevDefaultBeaconConfig
		err = undo()
		if err != nil {
			t.Error(err)
		}
		networkConfig = prevNetworkCfg
	})
}

func SetupForkEpochConfigForTest() {
	cfg := BeaconConfig().Copy()
	// original fork epoch
	// https://github.com/prysmaticlabs/prysm/blob/3413d05b3421c27579cf7a186cb9b142ffbb8346/config/params/mainnet_config.go#L25
	cfg.AltairForkEpoch = 72740
	cfg.BellatrixForkEpoch = 144896
	cfg.CapellaForkEpoch = 1904048
	cfg.InitializeForkSchedule()
	configs = newConfigset(cfg)
	OverrideBeaconConfig(cfg)
}
