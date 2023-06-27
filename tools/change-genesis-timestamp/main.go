package main

import (
	"flag"

	"github.com/prysmaticlabs/prysm/v4/encoding/ssz/detect"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
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

	// temp, err := detect.
	detector, err := detect.FromState(enc)
	if err != nil {
		panic(err)
	}

	// forkName := version.String(detector.Fork)

	switch fork := detector.Fork; fork {
	case version.Phase0:
		st := &ethpb.BeaconState{}
		err = st.UnmarshalSSZ(enc)
		if err != nil {
			panic(err)
		}
		st.GenesisTime = *timestamp
		enc, err = st.MarshalSSZ()
		if err != nil {
			panic(err)
		}

		err = file.WriteFile(*genesisPath, enc)
		if err != nil {
			panic(err)
		}
		return
	case version.Altair:
		st := &ethpb.BeaconStateAltair{}
		err = st.UnmarshalSSZ(enc)
		if err != nil {
			panic(err)
		}
		st.GenesisTime = *timestamp
		enc, err = st.MarshalSSZ()
		if err != nil {
			panic(err)
		}
		err = file.WriteFile(*genesisPath, enc)
		if err != nil {
			panic(err)
		}
		return
	case version.Bellatrix:
		st := &ethpb.BeaconStateBellatrix{}
		err = st.UnmarshalSSZ(enc)
		if err != nil {
			panic(err)
		}
		st.GenesisTime = *timestamp
		enc, err = st.MarshalSSZ()
		if err != nil {
			panic(err)
		}
		err = file.WriteFile(*genesisPath, enc)
		if err != nil {
			panic(err)
		}
		return
	case version.Capella:
		st := &ethpb.BeaconStateCapella{}
		err = st.UnmarshalSSZ(enc)
		if err != nil {
			panic(err)
		}
		st.GenesisTime = *timestamp
		enc, err = st.MarshalSSZ()
		if err != nil {
			panic(err)
		}
		err = file.WriteFile(*genesisPath, enc)
		if err != nil {
			panic(err)
		}
		return
	}

	// beaconState, err := detector.UnmarshalBeaconState(enc)
	// if err != nil {
	// 	panic(err)
	// }
	// s := beaconState.WriteOnlyBeaconState
	// err = s.SetGenesisTime(*timestamp)
	// // err = beaconState.WriteOnlyBeaconState().SetGenesisTime(*timestamp)
	// if err != nil {
	// 	panic(err)
	// }

	// // beaconState.GenesisTime = *timestamp

	// protoState := &ethpb.BeaconState{}
	// // if err := protoState.UnmarshalSSZ(enc); err != nil {
	// // 	panic(err)
	// // }
	// // protoState.GenesisTime = *timestamp

	// enc, err = beaconState.MarshalSSZ()
	// if err != nil {
	// 	panic(err)
	// }

	// err = file.WriteFile(*genesisPath, enc)
	// if err != nil {
	// 	panic(err)
	// }
}
