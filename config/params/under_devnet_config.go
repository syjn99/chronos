package params

import (
	"math"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
)

// UnderDevnetSpecConfig retrieves the mainnet undernet config used in spec tests.
func UnderDevnetSpecConfig() *BeaconChainConfig {
	underConfig := mainnetBeaconConfig.Copy()
	// Misc
	// underConfig.MaxCommitteesPerSlot = 4
	// underConfig.TargetCommitteeSize = 2
	// underConfig.MaxValidatorsPerCommittee = 2048
	// underConfig.MinPerEpochChurnLimit = 4
	// underConfig.ChurnLimitQuotient = 32
	// underConfig.ShuffleRoundCount = 10
	// underConfig.MinGenesisActiveValidatorCount = 64
	// underConfig.MinGenesisTime = 1685040000
	underConfig.GenesisDelay = 30 // 5 minutes
	// underConfig.TargetAggregatorsPerCommittee = 2

	// Gwei values
	underConfig.MinDepositAmount = 1e9
	underConfig.MaxEffectiveBalance = 32e9
	underConfig.EjectionBalance = 16e9
	underConfig.EffectiveBalanceIncrement = 1e9

	// Initial values
	// underConfig.BLSWithdrawalPrefixByte = byte(0)
	// underConfig.ETH1AddressWithdrawalPrefixByte = byte(1)

	// Time parameters
	// underConfig.SecondsPerSlot = 5
	// underConfig.MinAttestationInclusionDelay = 1
	// underConfig.SlotsPerEpoch = 8
	// underConfig.SqrRootSlotsPerEpoch = 2
	// underConfig.MinSeedLookahead = 1
	// underConfig.MaxSeedLookahead = 4
	// underConfig.EpochsPerEth1VotingPeriod = 4
	// underConfig.SlotsPerHistoricalRoot = 64
	// underConfig.MinValidatorWithdrawabilityDelay = 256
	// underConfig.ShardCommitteePeriod = 64
	// underConfig.MinEpochsToInactivityPenalty = 4
	// underConfig.Eth1FollowDistance = 16
	// underConfig.SecondsPerETH1Block = 5

	// State vector lengths
	// underConfig.EpochsPerHistoricalVector = 64
	// underConfig.EpochsPerSlashingsVector = 64
	// underConfig.HistoricalRootsLimit = 16777216
	// underConfig.ValidatorRegistryLimit = 1099511627776

	// Reward and penalty quotients
	// underConfig.BaseRewardFactor = 64
	// underConfig.WhistleBlowerRewardQuotient = 512
	// underConfig.ProposerRewardQuotient = 8
	// underConfig.InactivityPenaltyQuotient = 33554432
	// underConfig.MinSlashingPenaltyQuotient = 64
	// underConfig.ProportionalSlashingMultiplier = 2

	// // Max operations per block
	// underConfig.MaxProposerSlashings = 16
	// underConfig.MaxAttesterSlashings = 2
	// underConfig.MaxAttestations = 128
	// underConfig.MaxDeposits = 16
	// underConfig.MaxVoluntaryExits = 16
	// underConfig.MaxWithdrawalsPerPayload = 4
	// underConfig.MaxValidatorsPerWithdrawalsSweep = 16

	// Signature domains
	underConfig.DomainBeaconProposer = bytesutil.ToBytes4(bytesutil.Bytes4(0))
	underConfig.DomainBeaconAttester = bytesutil.ToBytes4(bytesutil.Bytes4(1))
	underConfig.DomainRandao = bytesutil.ToBytes4(bytesutil.Bytes4(2))
	underConfig.DomainDeposit = bytesutil.ToBytes4(hexutil.MustDecode("0x03000000"))
	underConfig.DomainVoluntaryExit = bytesutil.ToBytes4(bytesutil.Bytes4(4))
	// underConfig.GenesisForkVersion = []byte{0, 0, 0, 4}
	underConfig.GenesisForkVersion = (hexutil.MustDecode("0x20000089"))

	underConfig.DepositContractTreeDepth = 32
	underConfig.FarFutureEpoch = math.MaxUint64
	underConfig.FarFutureSlot = math.MaxUint64

	// New Altair params
	// underConfig.AltairForkVersion = []byte{1, 0, 0, 4} // Highest byte set to 0x01 to avoid collisions with mainnet versioning
	underConfig.AltairForkVersion = (hexutil.MustDecode("0x20000090"))
	underConfig.AltairForkEpoch = math.MaxUint64 - 1
	// underConfig.BellatrixForkVersion = []byte{2, 0, 0, 4}
	underConfig.BellatrixForkVersion = (hexutil.MustDecode("0x20000091"))
	underConfig.BellatrixForkEpoch = math.MaxUint64 - 1
	// underConfig.CapellaForkVersion = []byte{3, 0, 0, 4}
	underConfig.CapellaForkVersion = (hexutil.MustDecode("0x20000092"))
	underConfig.CapellaForkEpoch = math.MaxUint64 - 1

	// underConfig.SyncCommitteeSize = 32
	// underConfig.InactivityScoreBias = 4
	// underConfig.EpochsPerSyncCommitteePeriod = 8

	// Ethereum PoW parameters.
	underConfig.DepositChainID = 820   // Chain ID of eth1 under.
	underConfig.DepositNetworkID = 820 // Network ID of eth1 under.
	underConfig.DepositContractAddress = "000000000000000000000000000000000000beef"
	// 2**256-2**10 for fake minimal network
	underConfig.TerminalTotalDifficulty = "500"

	underConfig.ConfigName = UnderDevnetName
	underConfig.PresetBase = "under-devnet"

	underConfig.InitializeForkSchedule()
	return underConfig
}
