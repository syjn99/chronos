package main

import (
	"flag"

	"github.com/prysmaticlabs/prysm/v4/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
)

var (
	genesisPath = flag.String("genesis-state", "", "Path to genesis state file")
	timestamp   = flag.Uint64("timestamp", 0, "Timestamp to set in genesis state")
)

func main() {
	flag.Parse()
	enc, err := file.ReadFileAsBytes(*genesisPath)
	if err != nil {
		panic(err)
	}

	protoState := &ethpb.BeaconState{}
	if err := protoState.UnmarshalSSZ(enc); err != nil {
		panic(err)
	}
	protoState.GenesisTime = *timestamp

	enc, err = protoState.MarshalSSZ()
	if err != nil {
		panic(err)
	}

	err = file.WriteFile(*genesisPath, enc)
	if err != nil {
		panic(err)
	}
}
