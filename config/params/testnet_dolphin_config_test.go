package params_test

// TODO : Implement external/dolphin_testnet
//func TestDolphinConfigMatchesUpstreamYaml(t *testing.T) {
//	presetFPs := presetsFilePath(t, "mainnet")
//	mn, err := params.ByName(params.MainnetName)
//	require.NoError(t, err)
//	cfg := mn.Copy()
//	for _, fp := range presetFPs {
//		cfg, err = params.UnmarshalConfigFile(fp, cfg)
//		require.NoError(t, err)
//	}
//	fPath, err := bazel.Runfile("external/dolphin_testnet")
//	require.NoError(t, err)
//	configFP := path.Join(fPath, "custom_config_data", "config.yaml")
//	dcfg, err := params.UnmarshalConfigFile(configFP, nil)
//	require.NoError(t, err)
//	fields := fieldsFromYamls(t, append(presetFPs, configFP))
//	assertYamlFieldsMatch(t, "dolphin", fields, dcfg, params.DolphinConfig())
//}
