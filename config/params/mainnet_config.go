package params

import (
	"math"
	"time"

	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
)

// MainnetConfig returns the configuration to be used in the main network.
func MainnetConfig() *BeaconChainConfig {
	if mainnetBeaconConfig.ForkVersionSchedule == nil {
		mainnetBeaconConfig.InitializeForkSchedule()
	}
	mainnetBeaconConfig.InitializeEpochIssuance()
	mainnetBeaconConfig.InitializeDepositPlan()
	return mainnetBeaconConfig
}

const (
	// Genesis Fork Epoch for the mainnet config.
	genesisForkEpoch = 0
	// Altair Fork Epoch for mainnet config.
	mainnetAltairForkEpoch = 0
	// Bellatrix Fork Epoch for mainnet config.
	mainnetBellatrixForkEpoch = 0
	// Capella Fork Epoch for mainnet config.
	mainnetCapellaForkEpoch = 10
)

var mainnetNetworkConfig = &NetworkConfig{
	GossipMaxSize:                   1 << 20,      // 1 MiB
	GossipMaxSizeBellatrix:          10 * 1 << 20, // 10 MiB
	MaxChunkSize:                    1 << 20,      // 1 MiB
	MaxChunkSizeBellatrix:           10 * 1 << 20, // 10 MiB
	AttestationSubnetCount:          64,
	AttestationPropagationSlotRange: 32,
	MaxRequestBlocks:                1 << 10, // 1024
	TtfbTimeout:                     5 * time.Second,
	RespTimeout:                     10 * time.Second,
	MaximumGossipClockDisparity:     500 * time.Millisecond,
	MessageDomainInvalidSnappy:      [4]byte{00, 00, 00, 00},
	MessageDomainValidSnappy:        [4]byte{01, 00, 00, 00},
	ETH2Key:                         "over2",
	AttSubnetKey:                    "attnets",
	SyncCommsSubnetKey:              "syncnets",
	MinimumPeersInSubnetSearch:      20,
	ContractDeploymentBlock:         0, // Note: contract was deployed in genesis block.
	BootstrapNodes: []string{
		// Over Mainnet Bootnodes
		"enr:-LG4QPSpG5B2-RqEPFTIqioW9FPf3zbdihMzgcjq0EoDqrKWaHSNp9_cdJu1NqW801NQqbMIbu4ulMM6etZAUuzVfuOGAZAqEhQLh2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEIi9wdIRvdmVykIxhDBUAAAAS__________-Jc2VjcDI1NmsxoQKfJRWSiWQNWAC4yGT0aXhtsK9a2nHgWCtKZDSbmGAZD4N1ZHCCyyA", // Bootnode1
		"enr:-LG4QAeEVPoRW5d6wPm64dT0-q0dLxEHKN-GJteOZ9FzRuhbIODSAz6ri5aTh8h1RVBdtZCyzjQio9EapAgHo1Zx1nWGAZAqEhaqh2F0dG5ldHOIAAAAAAAAAACCaWSCdjSCaXCEIiABaYRvdmVykIxhDBUAAAAS__________-Jc2VjcDI1NmsxoQMutqRIiYjhL4_ZMDKyZBChUA7X2i7SmEx_Pqu7o8-LOIN1ZHCCyyA", // Bootnode2
	},
}

var mainnetBeaconConfig = &BeaconChainConfig{
	// Constants (Non-configurable)
	FarFutureEpoch:           math.MaxUint64,
	FarFutureSlot:            math.MaxUint64,
	BaseRewardsPerEpoch:      4,
	DepositContractTreeDepth: 32,
	GenesisDelay:             30, // 5 minutes

	// Misc constant.
	TargetCommitteeSize:               128,
	MaxValidatorsPerCommittee:         2048,
	MaxCommitteesPerSlot:              64,
	MinPerEpochChurnLimit:             4,
	ChurnLimitQuotient:                1 << 16,
	ChurnLimitBias:                    1,
	ShuffleRoundCount:                 90,
	MinGenesisActiveValidatorCount:    16384,
	MinGenesisTime:                    1718690400, // Jun 19, 2024, 00 AM UTC+9.
	TargetAggregatorsPerCommittee:     16,
	HysteresisQuotient:                4,
	HysteresisDownwardMultiplier:      1,
	HysteresisUpwardMultiplier:        5,
	DepositPlanEarlyEnd:               4,
	DepositPlanLaterEnd:               10,
	RewardFeedbackPrecision:           1000000000000,
	RewardFeedbackThresholdReciprocal: 10,
	TargetChangeRate:                  1500000,
	MaxBoostYield:                     10000000000,

	// Gwei value constants.
	MinDepositAmount:          1 * 1e9,
	MaxEffectiveBalance:       256 * 1e9,
	EjectionBalance:           128 * 1e9,
	EffectiveBalanceIncrement: 8 * 1e9,
	MaxTokenSupply:            1000000000 * 1e9,
	IssuancePerYear:           20000000 * 1e9,

	// Initial value constants.
	BLSWithdrawalPrefixByte:         byte(0),
	ETH1AddressWithdrawalPrefixByte: byte(1),
	ZeroHash:                        [32]byte{},

	// Time parameter constants.
	MinAttestationInclusionDelay:     1,
	SecondsPerSlot:                   12,
	SlotsPerEpoch:                    32,
	SqrRootSlotsPerEpoch:             5,
	EpochsPerYear:                    82125, // 365(Days)*24(Hours)*60(minutes)*60(Seconds)/12(Seconds per slot)/32(Slots per epoch)
	MinSeedLookahead:                 1,
	MaxSeedLookahead:                 4,
	EpochsPerEth1VotingPeriod:        64,
	SlotsPerHistoricalRoot:           8192,
	MinValidatorWithdrawabilityDelay: 256,
	ShardCommitteePeriod:             256,
	MinEpochsToInactivityPenalty:     4,
	Eth1FollowDistance:               1024,

	// Fork choice algorithm constants.
	ProposerScoreBoost:              40,
	ReorgWeightThreshold:            20,
	ReorgParentWeightThreshold:      160,
	ReorgMaxEpochsSinceFinalization: 2,
	IntervalsPerSlot:                3,

	// Ethereum PoW parameters.
	DepositChainID:         54176, // Chain ID of over mainnet.
	DepositNetworkID:       54176, // Network ID of over mainnet.
	DepositContractAddress: "000000000000000000000000000000000beac017",

	// Validator params.
	RandomSubnetsPerValidator:         1 << 0,
	EpochsPerRandomSubnetSubscription: 1 << 8,

	// While eth1 mainnet block times are closer to 13s, we must conform with other clients in
	// order to vote on the correct eth1 blocks.
	//
	// Additional context: https://github.com/ethereum/consensus-specs/issues/2132
	// Bug prompting this change: https://github.com/prysmaticlabs/prysm/issues/7856
	// Future optimization: https://github.com/prysmaticlabs/prysm/issues/7739
	SecondsPerETH1Block: 12,

	// State list length constants.
	EpochsPerHistoricalVector: 65536,
	EpochsPerSlashingsVector:  8192,
	HistoricalRootsLimit:      16777216,
	ValidatorRegistryLimit:    1099511627776,

	// Reward and penalty quotients constants.
	BaseRewardFactor:               64,
	WhistleBlowerRewardQuotient:    512,
	ProposerRewardQuotient:         8,
	InactivityPenaltyQuotient:      67108864,
	MinSlashingPenaltyQuotient:     128,
	ProportionalSlashingMultiplier: 1,

	// Max operations per block constants.
	MaxProposerSlashings:             16,
	MaxAttesterSlashings:             2,
	MaxAttestations:                  128,
	MaxDeposits:                      16,
	MaxVoluntaryExits:                16,
	MaxWithdrawalsPerPayload:         16,
	MaxBlsToExecutionChanges:         16,
	MaxValidatorsPerWithdrawalsSweep: 16384,

	// BLS domain values.
	DomainBeaconProposer:              bytesutil.Uint32ToBytes4(0x00000000),
	DomainBeaconAttester:              bytesutil.Uint32ToBytes4(0x01000000),
	DomainRandao:                      bytesutil.Uint32ToBytes4(0x02000000),
	DomainDeposit:                     bytesutil.Uint32ToBytes4(0x03000000),
	DomainVoluntaryExit:               bytesutil.Uint32ToBytes4(0x04000000),
	DomainSelectionProof:              bytesutil.Uint32ToBytes4(0x05000000),
	DomainAggregateAndProof:           bytesutil.Uint32ToBytes4(0x06000000),
	DomainSyncCommittee:               bytesutil.Uint32ToBytes4(0x07000000),
	DomainSyncCommitteeSelectionProof: bytesutil.Uint32ToBytes4(0x08000000),
	DomainContributionAndProof:        bytesutil.Uint32ToBytes4(0x09000000),
	DomainApplicationMask:             bytesutil.Uint32ToBytes4(0x00000001),
	DomainApplicationBuilder:          bytesutil.Uint32ToBytes4(0x00000001),
	DomainBLSToExecutionChange:        bytesutil.Uint32ToBytes4(0x0A000000),

	// Prysm constants.
	GweiPerEth:                     1000000000,
	BLSSecretKeyLength:             32,
	BLSPubkeyLength:                48,
	DefaultBufferSize:              10000,
	WithdrawalPrivkeyFileName:      "/shardwithdrawalkey",
	ValidatorPrivkeyFileName:       "/validatorprivatekey",
	RPCSyncCheck:                   1,
	EmptySignature:                 [96]byte{},
	DefaultPageSize:                250,
	MaxPeersToSync:                 15,
	SlotsPerArchivedPoint:          2048,
	GenesisCountdownInterval:       time.Minute,
	ConfigName:                     MainnetName,
	PresetBase:                     "mainnet",
	BeaconStateFieldCount:          24,
	BeaconStateAltairFieldCount:    28,
	BeaconStateBellatrixFieldCount: 29,
	BeaconStateCapellaFieldCount:   32,

	// Slasher related values.
	WeakSubjectivityPeriod:          54000,
	PruneSlasherStoragePeriod:       10,
	SlashingProtectionPruningEpochs: 512,

	// Weak subjectivity values.
	SafetyDecay: 10,

	// Fork related values.
	GenesisEpoch:         genesisForkEpoch,
	GenesisForkVersion:   []byte{0x00, 0x00, 0x00, 0x18},
	AltairForkVersion:    []byte{0x01, 0x00, 0x00, 0x18},
	AltairForkEpoch:      mainnetAltairForkEpoch,
	BellatrixForkVersion: []byte{0x02, 0x00, 0x00, 0x18},
	BellatrixForkEpoch:   mainnetBellatrixForkEpoch,
	CapellaForkVersion:   []byte{0x03, 0x00, 0x00, 0x18},
	CapellaForkEpoch:     mainnetCapellaForkEpoch,

	// New values introduced in Altair hard fork 1.
	// Participation flag indices.
	TimelySourceFlagIndex: 0,
	TimelyTargetFlagIndex: 1,
	TimelyHeadFlagIndex:   2,

	// Incentivization weight values.
	TimelySourceWeight: 12,
	TimelyTargetWeight: 24,
	TimelyHeadWeight:   12,
	SyncRewardWeight:   0,
	ProposerWeight:     8,
	LightLayerWeight:   8,
	WeightDenominator:  64,

	// Validator related values.
	TargetAggregatorsPerSyncSubcommittee: 16,
	SyncCommitteeSubnetCount:             4,

	// Misc values.
	SyncCommitteeSize:            512,
	InactivityScoreBias:          4,
	InactivityScoreRecoveryRate:  16,
	BailOutScoreBias:             1000000000000000,
	BailOutScoreThreshold:        1575000000000000000,
	EpochsPerSyncCommitteePeriod: 256,

	// Updated penalty values.
	InactivityPenaltyQuotientAltair:         3 * 1 << 24, // 50331648
	MinSlashingPenaltyQuotientAltair:        64,
	ProportionalSlashingMultiplierAltair:    2,
	MinSlashingPenaltyQuotientBellatrix:     32,
	ProportionalSlashingMultiplierBellatrix: 3,
	InactivityPenaltyQuotientBellatrix:      1 << 24,

	// Light client
	MinSyncCommitteeParticipants: 1,

	// Bellatrix
	TerminalBlockHashActivationEpoch: 18446744073709551615,
	TerminalBlockHash:                [32]byte{},
	TerminalTotalDifficulty:          "0",
	EthBurnAddressHex:                "0x0000000000000000000000000000000000000000",
	DefaultBuilderGasLimit:           uint64(30000000),

	// Mevboost circuit breaker
	MaxBuilderConsecutiveMissedSlots: 3,
	MaxBuilderEpochMissedSlots:       5,
	// Execution engine timeout value
	ExecutionEngineTimeoutValue: 8, // 8 seconds default based on: https://github.com/ethereum/execution-apis/blob/main/src/engine/specification.md#core
}

// MainnetTestConfig provides a version of the mainnet config that has a different name
// and a different fork choice schedule. This can be used in cases where we want to use config values
// that are consistent with mainnet, but won't conflict or cause the hard-coded genesis to be loaded.
func MainnetTestConfig() *BeaconChainConfig {
	mn := MainnetConfig().Copy()
	mn.ConfigName = MainnetTestName
	FillTestVersions(mn, 128)
	return mn
}

// FillTestVersions replaces the fork schedule in the given BeaconChainConfig with test values, using the given
// byte argument as the high byte (common across forks).
func FillTestVersions(c *BeaconChainConfig, b byte) {
	c.GenesisForkVersion = make([]byte, fieldparams.VersionLength)
	c.AltairForkVersion = make([]byte, fieldparams.VersionLength)
	c.BellatrixForkVersion = make([]byte, fieldparams.VersionLength)
	c.CapellaForkVersion = make([]byte, fieldparams.VersionLength)

	c.GenesisForkVersion[fieldparams.VersionLength-1] = b
	c.AltairForkVersion[fieldparams.VersionLength-1] = b
	c.BellatrixForkVersion[fieldparams.VersionLength-1] = b
	c.CapellaForkVersion[fieldparams.VersionLength-1] = b

	c.GenesisForkVersion[0] = 0
	c.AltairForkVersion[0] = 1
	c.BellatrixForkVersion[0] = 2
	c.CapellaForkVersion[0] = 3
}
