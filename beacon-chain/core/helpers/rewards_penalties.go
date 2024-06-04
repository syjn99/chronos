package helpers

import (
	"errors"
	"math"
	"math/big"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/cache"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/time"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	mathutil "github.com/prysmaticlabs/prysm/v4/math"
	"github.com/prysmaticlabs/prysm/v4/time/slots"
)

var balanceCache = cache.NewEffectiveBalanceCache()
var balanceWithQueueCache = cache.NewEffectiveBalanceCache()

// TotalBalance returns the total amount at stake in Gwei
// of input validators.
//
// Spec pseudocode definition:
//
//	def get_total_balance(state: BeaconState, indices: Set[ValidatorIndex]) -> Gwei:
//	 """
//	 Return the combined effective balance of the ``indices``.
//	 ``EFFECTIVE_BALANCE_INCREMENT`` Gwei minimum to avoid divisions by zero.
//	 Math safe up to ~10B ETH, afterwhich this overflows uint64.
//	 """
//	 return Gwei(max(EFFECTIVE_BALANCE_INCREMENT, sum([state.validators[index].effective_balance for index in indices])))
func TotalBalance(state state.ReadOnlyValidators, indices []primitives.ValidatorIndex) uint64 {
	total := uint64(0)

	for _, idx := range indices {
		val, err := state.ValidatorAtIndexReadOnly(idx)
		if err != nil {
			continue
		}
		total += val.EffectiveBalance()
	}

	// EFFECTIVE_BALANCE_INCREMENT is the lower bound for total balance.
	if total < params.BeaconConfig().EffectiveBalanceIncrement {
		return params.BeaconConfig().EffectiveBalanceIncrement
	}

	return total
}

// TotalActiveBalance returns the total amount at stake in Gwei
// of active validators.
//
// Spec pseudocode definition:
//
//	def get_total_active_balance(state: BeaconState) -> Gwei:
//	 """
//	 Return the combined effective balance of the active validators.
//	 Note: ``get_total_balance`` returns ``EFFECTIVE_BALANCE_INCREMENT`` Gwei minimum to avoid divisions by zero.
//	 """
//	 return get_total_balance(state, set(get_active_validator_indices(state, get_current_epoch(state))))
func TotalActiveBalance(s state.ReadOnlyBeaconState) (uint64, error) {
	bal, err := balanceCache.Get(s)
	switch {
	case err == nil:
		return bal, nil
	case errors.Is(err, cache.ErrNotFound):
		// Do nothing if we receive a not found error.
	default:
		// In the event, we encounter another error we return it.
		return 0, err
	}

	total := uint64(0)
	epoch := slots.ToEpoch(s.Slot())
	if err := s.ReadFromEveryValidator(func(idx int, val state.ReadOnlyValidator) error {
		if IsActiveValidatorUsingTrie(val, epoch) {
			total += val.EffectiveBalance()
		}
		return nil
	}); err != nil {
		return 0, err
	}

	// Spec defines `EffectiveBalanceIncrement` as min to avoid divisions by zero.
	total = mathutil.Max(params.BeaconConfig().EffectiveBalanceIncrement, total)
	if err := balanceCache.AddTotalEffectiveBalance(s, total); err != nil {
		return 0, err
	}

	return total, nil
}

// TotalBalanceWithQueue returns the total amount at stake in Gwei
// of active validators + pending validators - exiting validators.
func TotalBalanceWithQueue(s state.ReadOnlyBeaconState) (uint64, error) {
	bal, err := balanceWithQueueCache.Get(s)
	switch {
	case err == nil:
		return bal, nil
	case errors.Is(err, cache.ErrNotFound):
		// Do nothing if we receive a not found error.
	default:
		// In the event, we encounter another error we return it.
		return 0, err
	}

	total := uint64(0)
	totalExit := uint64(0)
	epoch := slots.ToEpoch(s.Slot())
	if err := s.ReadFromEveryValidator(func(idx int, val state.ReadOnlyValidator) error {
		if IsPendingValidatorUsingTrie(val, epoch) {
			total += val.EffectiveBalance()
		} else if IsActiveValidatorUsingTrie(val, epoch) {
			total += val.EffectiveBalance()
		}
		if IsExitingValidatorUsingTrie(val, epoch) {
			totalExit += val.EffectiveBalance()
		}
		return nil
	}); err != nil {
		return 0, err
	}
	if total > totalExit {
		total -= totalExit
	} else {
		total = 0
	}

	// Spec defines `EffectiveBalanceIncrement` as min to avoid divisions by zero.
	total = mathutil.Max(params.BeaconConfig().EffectiveBalanceIncrement, total)
	if err := balanceWithQueueCache.AddTotalEffectiveBalance(s, total); err != nil {
		return 0, err
	}

	return total, nil
}

// IncreaseBalance increases validator with the given 'index' balance by 'delta' in Gwei.
//
// Spec pseudocode definition:
//
//	def increase_balance(state: BeaconState, index: ValidatorIndex, delta: Gwei) -> None:
//	  """
//	  Increase the validator balance at index ``index`` by ``delta``.
//	  """
//	  state.balances[index] += delta
func IncreaseBalance(state state.BeaconState, idx primitives.ValidatorIndex, delta uint64) error {
	balAtIdx, err := state.BalanceAtIndex(idx)
	if err != nil {
		return err
	}
	newBal, err := IncreaseBalanceWithVal(balAtIdx, delta)
	if err != nil {
		return err
	}
	return state.UpdateBalancesAtIndex(idx, newBal)
}

// IncreaseBalanceWithVal increases validator with the given 'index' balance by 'delta' in Gwei.
// This method is flattened version of the spec method, taking in the raw balance and returning
// the post balance.
//
// Spec pseudocode definition:
//
//	def increase_balance(state: BeaconState, index: ValidatorIndex, delta: Gwei) -> None:
//	  """
//	  Increase the validator balance at index ``index`` by ``delta``.
//	  """
//	  state.balances[index] += delta
func IncreaseBalanceWithVal(currBalance, delta uint64) (uint64, error) {
	return mathutil.Add64(currBalance, delta)
}

// DecreaseBalance decreases validator with the given 'index' balance by 'delta' in Gwei.
//
// Spec pseudocode definition:
//
//	def decrease_balance(state: BeaconState, index: ValidatorIndex, delta: Gwei) -> None:
//	  """
//	  Decrease the validator balance at index ``index`` by ``delta``, with underflow protection.
//	  """
//	  state.balances[index] = 0 if delta > state.balances[index] else state.balances[index] - delta
func DecreaseBalance(state state.BeaconState, idx primitives.ValidatorIndex, delta uint64) error {
	balAtIdx, err := state.BalanceAtIndex(idx)
	if err != nil {
		return err
	}
	return state.UpdateBalancesAtIndex(idx, DecreaseBalanceWithVal(balAtIdx, delta))
}

// DecreaseBalanceWithVal decreases validator with the given 'index' balance by 'delta' in Gwei.
// This method is flattened version of the spec method, taking in the raw balance and returning
// the post balance.
//
// Spec pseudocode definition:
//
//	def decrease_balance(state: BeaconState, index: ValidatorIndex, delta: Gwei) -> None:
//	  """
//	  Decrease the validator balance at index ``index`` by ``delta``, with underflow protection.
//	  """
//	  state.balances[index] = 0 if delta > state.balances[index] else state.balances[index] - delta
func DecreaseBalanceWithVal(currBalance, delta uint64) uint64 {
	if delta > currBalance {
		return 0
	}
	return currBalance - delta
}

// IsInInactivityLeak returns true if the state is experiencing inactivity leak.
//
// Spec code:
// def is_in_inactivity_leak(state: BeaconState) -> bool:
//
//	return get_finality_delay(state) > MIN_EPOCHS_TO_INACTIVITY_PENALTY
func IsInInactivityLeak(prevEpoch, finalizedEpoch primitives.Epoch) bool {
	return FinalityDelay(prevEpoch, finalizedEpoch) > params.BeaconConfig().MinEpochsToInactivityPenalty
}

// FinalityDelay returns the finality delay using the beacon state.
//
// Spec code:
// def get_finality_delay(state: BeaconState) -> uint64:
//
//	return get_previous_epoch(state) - state.finalized_checkpoint.epoch
func FinalityDelay(prevEpoch, finalizedEpoch primitives.Epoch) primitives.Epoch {
	return prevEpoch - finalizedEpoch
}

// EpochIssuance returns the total amount of ETH(in Gwei) to be issued in the given epoch.
func EpochIssuance(epoch primitives.Epoch) uint64 {
	cfg := params.BeaconConfig()
	year := uint64(epoch.Div(cfg.EpochsPerYear))

	// After year 10, no more issuance left.
	if year > 9 {
		return 0
	}
	return cfg.EpochIssuance
}

// TargetDepositPlan returns the target deposit plan for the given epoch.
func TargetDepositPlan(epoch primitives.Epoch) uint64 {
	cfg := params.BeaconConfig()
	e := uint64(epoch)
	if e < cfg.EpochsPerYear*cfg.DepositPlanEarlyEnd {
		return cfg.DepositPlanEarlySlope*e + cfg.DepositPlanEarlyOffset
	} else if e < cfg.EpochsPerYear*cfg.DepositPlanLaterEnd {
		return cfg.DepositPlanLaterSlope*e + cfg.DepositPlanLaterOffset
	} else {
		return cfg.DepositPlanFinal
	}
}

func ProcessRewardfactorUpdate(state state.BeaconState) error {
	// update reward adjustment factor
	calculatedFactor, err := CalculateRewardAdjustmentFactor(state)
	if err != nil {
		return err
	}
	err = state.SetRewardAdjustmentFactor(calculatedFactor)
	if err != nil {
		return err
	}

	// update previous epoch reserve
	err = state.SetPreviousEpochReserve(state.CurrentEpochReserve())
	if err != nil {
		return err
	}

	return nil
}

// TotalRewardWithReserveUsage returns the total reward and total reserve usage in the given epoch.
// Reserve is always depleted in Over tokenomics.
func TotalRewardWithReserveUsage(s state.ReadOnlyBeaconState) (uint64, uint64) {
	epochIssuance := EpochIssuance(time.CurrentEpoch(s))
	feedbackBoost := EpochFeedbackBoost(s)

	return epochIssuance + feedbackBoost, feedbackBoost
}

// EpochFeedbackBoost returns the boost reward from feedback model in tokenomics.
func EpochFeedbackBoost(s state.ReadOnlyBeaconState) uint64 {
	cfg := params.BeaconConfig()
	rewardAdjustmentFactor := s.RewardAdjustmentFactor()
	feedbackBoost := cfg.MaxTokenSupply / cfg.RewardFeedbackPrecision * rewardAdjustmentFactor / cfg.EpochsPerYear

	if reserve := s.PreviousEpochReserve(); feedbackBoost > reserve {
		return reserve
	}
	return feedbackBoost
}

// CalcutateRewardAdjustmentFactor returns the adjustment factor for the next epoch.
// This is pseudo code from the spec.
//
// Spec code:
//
//	def process_validator_reward(epoch, parent, deposit, pending_deposit, exit_deposit):
//		deposit_delta = pending_deposit - exit_deposit
//		future_deposit = deposit + deposit_delta
//		target_deposit = TARGET_DEPOSIT_PLAN[epoch]
//
//		error_rate = abs(future_deposit - target_deposit) / target_deposit
//		mitigating_factor = max(1e-6, min(1.0, error_rate / threshold))
//		change_rate = target_change_rate * mitigating_factor
//
//		if future_deposit >= target_deposit:
//			bias_delta = -change_rate
//		else:
//			bias_delta =  change_rate
//
//		prev_bias = parent.bias
//		bias = prev_bias + bias_delta
//		bias = max(MIN_YIELD - TARGET_YIELD, bias)
//		bias = min(bias, INCENIVE_LIMIT * TARGET_YIELD)
//
//		return bias
func CalculateRewardAdjustmentFactor(state state.ReadOnlyBeaconState) (uint64, error) {
	cfg := params.BeaconConfig()
	futureDeposit, err := TotalBalanceWithQueue(state)
	if err != nil {
		return 0, err
	}
	targetDeposit := TargetDepositPlan(time.NextEpoch(state))

	// Using big integers for precise calculation
	if futureDeposit >= math.MaxInt64 || targetDeposit >= math.MaxInt64 {
		return 0, errors.New("deposit exceeds max int64, cannot calculate reward adjustment factor")
	}
	bigFutureDeposit := big.NewInt(int64(futureDeposit)) // lint:ignore uintcast -- changeRate will not exceed int64 because of total issuance.
	bigTargetDeposit := big.NewInt(int64(targetDeposit)) // lint:ignore uintcast -- changeRate will not exceed int64 because of total issuance.
	bigRewardPrecision := big.NewInt(int64(cfg.RewardFeedbackPrecision))

	// Calculate the gap and error rate to make mitigating factor
	gap := new(big.Int).Abs(new(big.Int).Sub(bigFutureDeposit, bigTargetDeposit))
	numerator := new(big.Int).Mul(gap, bigRewardPrecision)
	errRate := new(big.Int).Div(numerator, bigTargetDeposit)
	mitigatingFactor := big.NewInt(int64(mathutil.Max(1000000, mathutil.Min(cfg.RewardFeedbackPrecision, errRate.Uint64()*cfg.RewardFeedbackThresholdReciprocal)))) // lint:ignore uintcast -- errRate will not exceed int64 because of truncation.

	// Calculate the change rate
	targetChangeRate := big.NewInt(int64(cfg.TargetChangeRate))
	changeRate := new(big.Int).Div(new(big.Int).Mul(targetChangeRate, mitigatingFactor), bigRewardPrecision).Uint64()
	if changeRate >= math.MaxUint64 {
		return 0, errors.New("changeRate exceeds max uint64, cannot calculate reward adjustment factor")
	}

	bias := state.RewardAdjustmentFactor()
	if futureDeposit >= targetDeposit {
		if bias <= changeRate {
			bias = 0
		} else {
			bias -= changeRate
		}
	} else {
		bias += changeRate
	}

	return TruncateRewardAdjustmentFactor(bias), nil
}

// TruncateRewardAdjustmentFactor truncates the given bias to the min and max bounds.
// The min and max bounds are defined in the spec.
func TruncateRewardAdjustmentFactor(bias uint64) uint64 {
	cfg := params.BeaconConfig()
	if 0 > bias {
		bias = 0
	}
	if cfg.MaxBoostYield < bias {
		bias = cfg.MaxBoostYield
	}
	return bias
}

// DecreaseCurrentReserve reduces the current reserve by the given amount.
// If the reserve is less than the given amount, it sets the reserve to 0.
func DecreaseCurrentReserve(state state.BeaconState, sub uint64) error {
	currentReserve := state.CurrentEpochReserve()
	if currentReserve < sub {
		err := state.SetCurrentEpochReserve(0)
		if err != nil {
			return err
		}
	} else {
		err := state.SetCurrentEpochReserve(currentReserve - sub)
		if err != nil {
			return err
		}
	}
	return nil
}
