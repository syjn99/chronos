package features

import (
	"time"

	"github.com/urfave/cli/v2"
)

var (
	// DolphinTestnet flag for the multiclient over consensus testnet.
	DolphinTestnet = &cli.BoolFlag{
		Name:  "dolphin",
		Usage: "Run Chronos configured for the Dolphin test network",
	}
	// Mainnet flag for easier tooling, no-op
	Mainnet = &cli.BoolFlag{
		Value: true,
		Name:  "mainnet",
		Usage: "Run on Over Protocol Beacon Chain Main Net. This is the default and can be omitted.",
	}
	devModeFlag = &cli.BoolFlag{
		Name:  "dev",
		Usage: "Enable experimental features still in development. These features may not be stable.",
	}
	writeSSZStateTransitionsFlag = &cli.BoolFlag{
		Name:  "interop-write-ssz-state-transitions",
		Usage: "Write ssz states to disk after attempted state transition",
	}
	enableExternalSlasherProtectionFlag = &cli.BoolFlag{
		Name: "enable-external-slasher-protection",
		Usage: "Enables the validator to connect to a beacon node using the --slasher flag" +
			"for remote slashing protection",
	}
	disableGRPCConnectionLogging = &cli.BoolFlag{
		Name:  "disable-grpc-connection-logging",
		Usage: "Disables displaying logs for newly connected grpc clients",
	}
	disableReorgLateBlocks = &cli.BoolFlag{
		Name:  "disable-reorg-late-blocks",
		Usage: "Disables reorgs of late blocks",
	}
	disablePeerScorer = &cli.BoolFlag{
		Name:  "disable-peer-scorer",
		Usage: "(Danger): Disables P2P peer scorer. Do NOT use this in production!",
	}
	disableCheckBadPeer = &cli.BoolFlag{
		Name:  "disable-check-bad-peer",
		Usage: "(Danger): Disables checking if a peer is bad. Do NOT use this in production!",
	}
	// writeWalletPasswordOnWebOnboarding = &cli.BoolFlag{
	// 	Name: "write-wallet-password-on-web-onboarding",
	// 	Usage: "(Danger): Writes the wallet password to the wallet directory on completing Prysm web onboarding. " +
	// 		"We recommend against this flag unless you are an advanced user.",
	// }
	aggregateFirstInterval = &cli.DurationFlag{
		Name:   "aggregate-first-interval",
		Usage:  "(Advanced): Specifies the first interval in which attestations are aggregated in the slot (typically unnaggregated attestations are aggregated in this interval)",
		Value:  6500 * time.Millisecond,
		Hidden: true,
	}
	aggregateSecondInterval = &cli.DurationFlag{
		Name:   "aggregate-second-interval",
		Usage:  "(Advanced): Specifies the second interval in which attestations are aggregated in the slot",
		Value:  9500 * time.Millisecond,
		Hidden: true,
	}
	aggregateThirdInterval = &cli.DurationFlag{
		Name:   "aggregate-third-interval",
		Usage:  "(Advanced): Specifies the third interval in which attestations are aggregated in the slot",
		Value:  11800 * time.Millisecond,
		Hidden: true,
	}
	dynamicKeyReloadDebounceInterval = &cli.DurationFlag{
		Name: "dynamic-key-reload-debounce-interval",
		Usage: "(Advanced): Specifies the time duration the validator waits to reload new keys if they have " +
			"changed on disk. Default 1s, can be any type of duration such as 1.5s, 1000ms, 1m.",
		Value: time.Second,
	}
	disableBroadcastSlashingFlag = &cli.BoolFlag{
		Name:  "disable-broadcast-slashings",
		Usage: "Disables broadcasting slashings submitted to the beacon node.",
	}
	attestTimely = &cli.BoolFlag{
		Name:  "attest-timely",
		Usage: "Fixes validator can attest timely after current block processes. See #8185 for more details",
	}
	enableSlasherFlag = &cli.BoolFlag{
		Name:  "slasher",
		Usage: "Enables a slasher in the beacon node for detecting slashable offenses",
	}
	enableSlashingProtectionPruning = &cli.BoolFlag{
		Name:  "enable-slashing-protection-history-pruning",
		Usage: "Enables the pruning of the validator client's slashing protection database",
	}
	enableDoppelGangerProtection = &cli.BoolFlag{
		Name: "enable-doppelganger",
		Usage: "Enables the validator to perform a doppelganger check. (Warning): This is not " +
			"a foolproof method to find duplicate instances in the network. Your validator will still be" +
			" vulnerable if it is being run in unsafe configurations.",
	}
	disableStakinContractCheck = &cli.BoolFlag{
		Name:  "disable-staking-contract-check",
		Usage: "Disables checking of staking contract deposits when proposing blocks, useful for devnets",
	}
	enableHistoricalSpaceRepresentation = &cli.BoolFlag{
		Name: "enable-historical-state-representation",
		Usage: "Enables the beacon chain to save historical states in a space efficient manner." +
			" (Warning): Once enabled, this feature migrates your database in to a new schema and " +
			"there is no going back. At worst, your entire database might get corrupted.",
	}
	enableStartupOptimistic = &cli.BoolFlag{
		Name:   "startup-optimistic",
		Usage:  "Treats every block as optimistically synced at launch. Use with caution",
		Value:  false,
		Hidden: true,
	}
	enableFullSSZDataLogging = &cli.BoolFlag{
		Name:  "enable-full-ssz-data-logging",
		Usage: "Enables displaying logs for full ssz data on rejected gossip messages",
	}
	SaveFullExecutionPayloads = &cli.BoolFlag{
		Name:  "save-full-execution-payloads",
		Usage: "Saves beacon blocks with full execution payloads instead of execution payload headers in the database",
	}
	EnableBeaconRESTApi = &cli.BoolFlag{
		Name:  "enable-beacon-rest-api",
		Usage: "Experimental enable of the beacon REST API when querying a beacon node",
	}
	enableVerboseSigVerification = &cli.BoolFlag{
		Name:  "enable-verbose-sig-verification",
		Usage: "Enables identifying invalid signatures if batch verification fails when processing block",
	}
	enableOptionalEngineMethods = &cli.BoolFlag{
		Name:  "enable-optional-engine-methods",
		Usage: "Enables the optional engine methods",
	}
	prepareAllPayloads = &cli.BoolFlag{
		Name:  "prepare-all-payloads",
		Usage: "Informs the engine to prepare all local payloads. Useful for relayers and builders",
	}
	disableBuildBlockParallel = &cli.BoolFlag{
		Name:  "disable-build-block-parallel",
		Usage: "Disables building a beacon block in parallel for consensus and execution",
	}
	disableResourceManager = &cli.BoolFlag{
		Name:  "disable-resource-manager",
		Usage: "Disables running the libp2p resource manager",
	}

	// DisableRegistrationCache a flag for disabling the validator registration cache and use db instead.
	DisableRegistrationCache = &cli.BoolFlag{
		Name:  "diable-registration-cache",
		Usage: "A temporary flag for disabling the validator registration cache instead of using the db. note: registrations do not clear on restart while using the db",
	}

	aggregateParallel = &cli.BoolFlag{
		Name:  "aggregate-parallel",
		Usage: "Enables parallel aggregation of attestations",
	}
)

// devModeFlags holds list of flags that are set when development mode is on.
var devModeFlags = []cli.Flag{
	enableVerboseSigVerification,
	enableOptionalEngineMethods,
}

// ValidatorFlags contains a list of all the feature flags that apply to the validator client.
var ValidatorFlags = append(deprecatedFlags, []cli.Flag{
	// writeWalletPasswordOnWebOnboarding,
	enableExternalSlasherProtectionFlag,
	DolphinTestnet,
	Mainnet,
	dynamicKeyReloadDebounceInterval,
	attestTimely,
	enableSlashingProtectionPruning,
	enableDoppelGangerProtection,
	EnableBeaconRESTApi,
}...)

// E2EValidatorFlags contains a list of the validator feature flags to be tested in E2E.
var E2EValidatorFlags = []string{
	"--enable-doppelganger",
}

// BeaconChainFlags contains a list of all the feature flags that apply to the beacon-chain client.
var BeaconChainFlags = append(deprecatedBeaconFlags, append(deprecatedFlags, []cli.Flag{
	devModeFlag,
	writeSSZStateTransitionsFlag,
	disableGRPCConnectionLogging,
	DolphinTestnet,
	Mainnet,
	disablePeerScorer,
	disableCheckBadPeer,
	disableBroadcastSlashingFlag,
	enableSlasherFlag,
	enableHistoricalSpaceRepresentation,
	disableStakinContractCheck,
	disableReorgLateBlocks,
	SaveFullExecutionPayloads,
	enableStartupOptimistic,
	enableFullSSZDataLogging,
	enableVerboseSigVerification,
	enableOptionalEngineMethods,
	prepareAllPayloads,
	disableBuildBlockParallel,
	aggregateFirstInterval,
	aggregateSecondInterval,
	aggregateThirdInterval,
	disableResourceManager,
	DisableRegistrationCache,
	aggregateParallel,
}...)...)

// E2EBeaconChainFlags contains a list of the beacon chain feature flags to be tested in E2E.
var E2EBeaconChainFlags = []string{
	"--dev",
}

// NetworkFlags contains a list of network flags.
var NetworkFlags = []cli.Flag{
	Mainnet,
	DolphinTestnet,
}
