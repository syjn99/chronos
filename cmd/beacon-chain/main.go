// Package beacon-chain defines the entire runtime of an Ethereum beacon node.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	runtimeDebug "runtime/debug"

	golog "github.com/ipfs/go-log/v2"
	joonix "github.com/joonix/log"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/builder"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/node"
	"github.com/prysmaticlabs/prysm/v4/cmd"
	blockchaincmd "github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/blockchain"
	dbcommands "github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/db"
	"github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/flags"
	jwtcommands "github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/jwt"
	"github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/sync/checkpoint"
	"github.com/prysmaticlabs/prysm/v4/cmd/beacon-chain/sync/genesis"
	"github.com/prysmaticlabs/prysm/v4/config/features"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	"github.com/prysmaticlabs/prysm/v4/io/logs"
	"github.com/prysmaticlabs/prysm/v4/monitoring/journald"
	"github.com/prysmaticlabs/prysm/v4/runtime/debug"
	"github.com/prysmaticlabs/prysm/v4/runtime/fdlimits"
	prefixed "github.com/prysmaticlabs/prysm/v4/runtime/logging/logrus-prefixed-formatter"
	_ "github.com/prysmaticlabs/prysm/v4/runtime/maxprocs"
	"github.com/prysmaticlabs/prysm/v4/runtime/tos"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var appFlags = []cli.Flag{
	flags.DepositContractFlag,
	flags.ExecutionEngineEndpoint,
	flags.ExecutionEngineHeaders,
	flags.ExecutionJWTSecretFlag,
	flags.RPCHost,
	flags.RPCPort,
	flags.CertFlag,
	flags.KeyFlag,
	flags.HTTPModules,
	flags.DisableGRPCGateway,
	flags.GRPCGatewayHost,
	flags.GRPCGatewayPort,
	flags.GPRCGatewayCorsDomain,
	flags.MinSyncPeers,
	flags.ContractDeploymentBlock,
	flags.SetGCPercent,
	flags.BlockBatchLimit,
	flags.BlockBatchLimitBurstFactor,
	flags.InteropMockEth1DataVotesFlag,
	flags.InteropNumValidatorsFlag,
	flags.InteropGenesisTimeFlag,
	flags.SlotsPerArchivedPoint,
	flags.EnableDebugRPCEndpoints,
	flags.EnableOverNodeRPCEndpoints,
	flags.SubscribeToAllSubnets,
	flags.HistoricalSlasherNode,
	flags.ChainID,
	flags.NetworkID,
	flags.WeakSubjectivityCheckpoint,
	flags.Eth1HeaderReqLimit,
	flags.MinPeersPerSubnet,
	flags.SuggestedFeeRecipient,
	flags.TerminalTotalDifficultyOverride,
	flags.TerminalBlockHashOverride,
	flags.TerminalBlockHashActivationEpochOverride,
	// flags.MevRelayEndpoint, // Temporarily deactivate for operational verification.
	flags.MaxBuilderEpochMissedSlots,
	flags.MaxBuilderConsecutiveMissedSlots,
	flags.EngineEndpointTimeoutSeconds,
	cmd.BackupWebhookOutputDir,
	cmd.MinimalConfigFlag,
	cmd.E2EConfigFlag,
	cmd.RPCMaxPageSizeFlag,
	cmd.BootstrapNode,
	cmd.NoDiscovery,
	cmd.StaticPeers,
	cmd.RelayNode,
	cmd.P2PUDPPort,
	cmd.P2PTCPPort,
	cmd.P2PIP,
	cmd.P2PHost,
	cmd.P2PHostDNS,
	cmd.P2PMaxPeers,
	cmd.P2PPrivKey,
	cmd.P2PStaticID,
	cmd.P2PMetadata,
	cmd.P2PAllowList,
	cmd.P2PDenyList,
	cmd.P2PColocationWhitelistFlag,
	cmd.P2PColocationLimitFlag,
	cmd.P2PIpTrackerBanTimeFlag,
	cmd.DataDirFlag,
	cmd.VerbosityFlag,
	cmd.EnableTracingFlag,
	cmd.TracingProcessNameFlag,
	cmd.TracingEndpointFlag,
	cmd.TraceSampleFractionFlag,
	cmd.MonitoringHostFlag,
	flags.MonitoringPortFlag,
	cmd.DisableMonitoringFlag,
	cmd.ClearDB,
	cmd.ForceClearDB,
	cmd.LogFormat,
	cmd.MaxGoroutines,
	debug.PProfFlag,
	debug.PProfAddrFlag,
	debug.PProfPortFlag,
	debug.MemProfileRateFlag,
	debug.CPUProfileFlag,
	debug.TraceFlag,
	debug.BlockProfileRateFlag,
	debug.MutexProfileFractionFlag,
	cmd.LogFileName,
	cmd.LogRotateFlag,
	cmd.LogMaxSizeMBsFlag,
	cmd.LogMaxBackupsFlag,
	cmd.LogMaxAgeFlag,
	cmd.LogCompressFlag,
	cmd.EnableUPnPFlag,
	cmd.ConfigFileFlag,
	cmd.ChainConfigFileFlag,
	cmd.GrpcMaxCallRecvMsgSizeFlag,
	cmd.AcceptTosFlag,
	cmd.RestoreSourceFileFlag,
	cmd.RestoreTargetDirFlag,
	cmd.ValidatorMonitorIndicesFlag,
	cmd.ApiTimeoutFlag,
	checkpoint.BlockPath,
	checkpoint.StatePath,
	checkpoint.RemoteURL,
	genesis.StatePath,
	genesis.BeaconAPIURL,
	flags.SlasherDirFlag,
}

func init() {
	appFlags = cmd.WrapFlags(append(appFlags, features.BeaconChainFlags...))
}

func main() {
	app := cli.App{}
	app.Name = "beacon-chain"
	app.Usage = "this is a beacon chain implementation for Over Protocol"
	app.Action = func(ctx *cli.Context) error {
		if err := startNode(ctx); err != nil {
			return cli.Exit(err.Error(), 1)
		}
		return nil
	}
	app.Version = version.Version()
	app.Commands = []*cli.Command{
		dbcommands.Commands,
		jwtcommands.Commands,
	}

	app.Flags = appFlags
	app.Before = func(ctx *cli.Context) error {
		// Load flags from config file, if specified.
		if err := cmd.LoadFlagsFromConfig(ctx, app.Flags); err != nil {
			return err
		}

		if err := cmd.ExpandSingleEndpointIfFile(ctx, flags.ExecutionEngineEndpoint); err != nil {
			return err
		}
		if ctx.IsSet(flags.SetGCPercent.Name) {
			runtimeDebug.SetGCPercent(ctx.Int(flags.SetGCPercent.Name))
		}
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		if err := fdlimits.SetMaxFdLimits(); err != nil {
			return err
		}
		return cmd.ValidateNoArgs(ctx)
	}

	defer func() {
		if x := recover(); x != nil {
			log.Errorf("Runtime panic: %v\n%v", x, string(runtimeDebug.Stack()))
			panic(x)
		}
	}()

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

func startNode(ctx *cli.Context) error {
	// Fix data dir for Windows users.
	outdatedDataDir := filepath.Join(file.HomeDir(), "AppData", "Roaming", "Eth2")
	currentDataDir := ctx.String(cmd.DataDirFlag.Name)
	if err := cmd.FixDefaultDataDir(outdatedDataDir, currentDataDir); err != nil {
		return err
	}

	// verify if ToS accepted
	if err := tos.VerifyTosAcceptedOrPrompt(ctx); err != nil {
		return err
	}

	verbosity := ctx.String(cmd.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)
	// Set libp2p logger to only panic logs for the info level.
	golog.SetAllLoggers(golog.LevelPanic)

	if level == logrus.DebugLevel {
		// Set libp2p logger to error logs for the debug level.
		golog.SetAllLoggers(golog.LevelError)
	}
	if level == logrus.TraceLevel {
		// libp2p specific logging.
		golog.SetAllLoggers(golog.LevelDebug)
		// Geth specific logging.
		// glogger := gethlog.NewGlogHandler(gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true)))
		// glogger.Verbosity(gethlog.LvlTrace)
		// gethlog.Root().SetHandler(glogger)
	}

	format := ctx.String(cmd.LogFormat.Name)
	switch format {
	case "text":
		formatter := new(prefixed.TextFormatter)
		formatter.TimestampFormat = "2006-01-02 15:04:05"
		formatter.FullTimestamp = true
		// If persistent log files are written - we disable the log messages coloring because
		// the colors are ANSI codes and seen as gibberish in the log files.
		formatter.DisableColors = ctx.String(cmd.LogFileName.Name) != ""
		logrus.SetFormatter(formatter)
	case "fluentd":
		f := joonix.NewFormatter()
		if err := joonix.DisableTimestampFormat(f); err != nil {
			panic(err)
		}
		logrus.SetFormatter(f)
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "journald":
		if err := journald.Enable(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown log format %s", format)
	}

	// Set log rotation if log file name is specified.
	logFileName := ctx.String(cmd.LogFileName.Name)
	rotation := ctx.Bool(cmd.LogRotateFlag.Name)
	if logFileName != "" {
		logrus.SetOutput(os.Stdout) // Set default output to stdout for log separation
		if rotation {
			lumberjackLogger := &lumberjack.Logger{
				Filename:   logFileName,
				MaxSize:    ctx.Int(cmd.LogMaxSizeMBsFlag.Name), // MB
				MaxBackups: ctx.Int(cmd.LogMaxBackupsFlag.Name),
				MaxAge:     ctx.Int(cmd.LogMaxAgeFlag.Name), // days
				Compress:   ctx.Bool(cmd.LogCompressFlag.Name),
			}
			logs.AddLogWriter(io.MultiWriter(logrus.StandardLogger().Out, lumberjackLogger))
			log.Info("Log rotation activated. path=", logFileName, ", MaxSize=", lumberjackLogger.MaxSize, ", MaxBackups=", lumberjackLogger.MaxBackups, ", MaxAge=", lumberjackLogger.MaxAge, ", Compress=", lumberjackLogger.Compress)
		} else if err := logs.ConfigurePersistentLogging(logFileName); err != nil {
			log.WithError(err).Error("Failed to configuring logging to disk.")
		}
	}

	blockchainFlagOpts, err := blockchaincmd.FlagOptions(ctx)
	if err != nil {
		return err
	}
	executionFlagOpts, err := execution.FlagOptions(ctx)
	if err != nil {
		return err
	}
	builderFlagOpts, err := builder.FlagOptions(ctx)
	if err != nil {
		return err
	}
	opts := []node.Option{
		node.WithBlockchainFlagOptions(blockchainFlagOpts),
		node.WithExecutionChainOptions(executionFlagOpts),
		node.WithBuilderFlagOptions(builderFlagOpts),
	}

	optFuncs := []func(*cli.Context) (node.Option, error){
		genesis.BeaconNodeOptions,
		checkpoint.BeaconNodeOptions,
	}
	for _, of := range optFuncs {
		ofo, err := of(ctx)
		if err != nil {
			return err
		}
		if ofo != nil {
			opts = append(opts, ofo)
		}
	}

	beacon, err := node.New(ctx, opts...)
	if err != nil {
		return err
	}
	beacon.Start()
	return nil
}
