package main

import (
	_ "embed"
	"flag"

	"github.com/golang/snappy"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

var (
	inputPath  = flag.String("input-path", "", "The input file path (/path/to/old.ssz)")
	outputPath = flag.String("output-path", "", "The output file path (/path/to/new.ssz)")
	oldSSZ     []byte
	isSnappy   bool
)

func main() {
	flag.Parse()
	enc, err := file.ReadFileAsBytes(*inputPath)
	if err != nil {
		panic(err)
	}
	// try decode snappy
	isSnappy = true
	oldSSZ, err = snappy.Decode(nil, enc)
	if err != nil {
		isSnappy = false
		oldSSZ = enc
	}

	oldSt := &ethpb.BeaconStateOld{}
	err = oldSt.UnmarshalSSZ(oldSSZ)
	if err != nil {
		panic(err)
	}
	newSt := &ethpb.BeaconState{}
	// construct new state from old state
	newSt.GenesisTime = oldSt.GenesisTime
	newSt.GenesisValidatorsRoot = oldSt.GenesisValidatorsRoot
	newSt.Slot = oldSt.Slot
	newSt.Fork = oldSt.Fork
	newSt.LatestBlockHeader = oldSt.LatestBlockHeader
	newSt.BlockRoots = oldSt.BlockRoots
	newSt.StateRoots = oldSt.StateRoots
	newSt.HistoricalRoots = oldSt.HistoricalRoots
	newSt.Eth1Data = oldSt.Eth1Data
	newSt.Eth1DataVotes = oldSt.Eth1DataVotes
	newSt.Eth1DepositIndex = oldSt.Eth1DepositIndex
	newSt.Validators = oldSt.Validators
	newSt.Balances = oldSt.Balances
	newSt.RandaoMixes = oldSt.RandaoMixes
	newSt.Slashings = oldSt.Slashings
	newSt.PreviousEpochAttestations = oldSt.PreviousEpochAttestations
	newSt.CurrentEpochAttestations = oldSt.CurrentEpochAttestations
	newSt.JustificationBits = oldSt.JustificationBits
	newSt.PreviousJustifiedCheckpoint = oldSt.PreviousJustifiedCheckpoint
	newSt.CurrentJustifiedCheckpoint = oldSt.CurrentJustifiedCheckpoint
	newSt.FinalizedCheckpoint = oldSt.FinalizedCheckpoint

	// Fill: RewardAdjustmentFactor, PreviousEpochReserve, CurrentEpochReserve
	newSt.RewardAdjustmentFactor = 0
	newSt.PreviousEpochReserve = 0
	newSt.CurrentEpochReserve = 0
	nb, err := newSt.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	// encode snappy
	if isSnappy {
		nb = snappy.Encode(nil, nb)
	}
	err = file.WriteFile(*outputPath, nb)
	if err != nil {
		panic(err)
	}
}
