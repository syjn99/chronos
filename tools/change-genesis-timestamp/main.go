package main

import (
	"flag"

	"github.com/prysmaticlabs/prysm/v4/encoding/ssz/detect"
	"github.com/prysmaticlabs/prysm/v4/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
	prysmTime "github.com/prysmaticlabs/prysm/v4/time"
)

var (
	genesisPath = flag.String("genesis-state", "", "Path to genesis state file")
	timestamp   = flag.Uint64("timestamp", 0, "Timestamp to set in genesis state")
)

func main() {
	flag.Parse()

	// if timestamp is not set, use current time + 60s
	if *timestamp == 0 {
		*timestamp = uint64(prysmTime.Now().Unix() + 60)
	}

	enc, err := file.ReadFileAsBytes(*genesisPath)
	if err != nil {
		panic(err)
	}

	detector, err := detect.FromState(enc)
	if err != nil {
		panic(err)
	}

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
}
