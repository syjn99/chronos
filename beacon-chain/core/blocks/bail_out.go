package blocks

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	v "github.com/prysmaticlabs/prysm/v4/beacon-chain/core/validators"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
	"github.com/sirupsen/logrus"
)

// ValidatorCannotBailoutYetMsg defines a message saying that a validator has already exited.
var ValidatorCannotBailoutYetMsg = "validator bailout score is below bail out threshold"

// ProcessBailOuts is one of the operations performed
// on each processed beacon block to determine which validators
// should exit the state's validator registry.
//
// Spec pseudocode definition:
//
//	def process_voluntary_exit(state: BeaconState, signed_voluntary_exit: SignedVoluntaryExit) -> None:
//	 voluntary_exit = signed_voluntary_exit.message
//	 validator = state.validators[voluntary_exit.validator_index]
//	 # Verify the validator is active
//	 assert is_active_validator(validator, get_current_epoch(state))
//	 # Verify exit has not been initiated
//	 assert validator.exit_epoch == FAR_FUTURE_EPOCH
//	 # Exits must specify an epoch when they become valid; they are not valid before then
//	 assert get_current_epoch(state) >= voluntary_exit.epoch
//	 # Verify the validator has been active long enough
//	 assert get_current_epoch(state) >= validator.activation_epoch + SHARD_COMMITTEE_PERIOD
//	 # Verify signature
//	 domain = get_domain(state, DOMAIN_VOLUNTARY_EXIT, voluntary_exit.epoch)
//	 signing_root = compute_signing_root(voluntary_exit, domain)
//	 assert bls.Verify(validator.pubkey, signing_root, signed_voluntary_exit.signature)
//	 # Initiate exit
//	 initiate_validator_exit(state, voluntary_exit.validator_index)
func ProcessBailOuts(
	ctx context.Context,
	beaconState state.BeaconState,
	exits []*ethpb.BailOut,
) (state.BeaconState, error) {
	for idx, exit := range exits {
		if exit == nil {
			return nil, errors.New("nil bail out in block body")
		}
		var err error
		valIdx := exit.ValidatorIndex
		if err = VerifyBailOut(beaconState, valIdx); err != nil {
			return nil, errors.Wrapf(err, "could not verify exit %d", idx)
		}
		beaconState, err = v.InitiateValidatorExit(ctx, beaconState, valIdx, true)
		if err != nil {
			return nil, err
		}
	}
	return beaconState, nil
}

// VerifyBailOut implements the validation for bail outs.
// It checks if validator is active, has not yet submitted an exit, is eligible to exit(above bail out score threshold)
func VerifyBailOut(
	state state.ReadOnlyBeaconState,
	valIdx primitives.ValidatorIndex,
) error {
	currentEpoch := slots.ToEpoch(state.Slot())
	// Verify validator index exists and get the validator.
	validator, err := state.ValidatorAtIndexReadOnly(valIdx)
	if err != nil {
		logrus.WithError(err).Warningf("could not get validator at index %d", valIdx)
		return err
	}
	// Verify the validator is active.
	if !helpers.IsActiveValidatorUsingTrie(validator, currentEpoch) {
		return errors.New("non-active validator cannot exit")
	}
	// Verify the validator has not yet submitted an exit.
	if validator.ExitEpoch() != params.BeaconConfig().FarFutureEpoch {
		return fmt.Errorf("validator with index %d %s: %v", valIdx, ValidatorAlreadyExitedMsg, validator.ExitEpoch())
	}
	// Verify the validator is eligible to exit.
	bailoutScores, err := state.BailOutScores()
	if err != nil {
		return err
	}
	if bailoutScores[valIdx] < params.BeaconConfig().BailOutScoreThreshold {
		return fmt.Errorf("validator with index %d %s: %v", valIdx, ValidatorCannotBailoutYetMsg, bailoutScores[valIdx])
	}

	return nil
}
