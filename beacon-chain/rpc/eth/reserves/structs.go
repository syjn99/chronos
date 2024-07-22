package reserves

type GetReservesResponse struct {
	Data                *Reserves `json:"data"`
	ExecutionOptimistic bool      `json:"execution_optimistic"`
	Finalized           bool      `json:"finalized"`
}

type Reserves struct {
	RewardAdjustmentFactor string `json:"reward_adjustment_factor"`
	PreviousEpochReserve   string `json:"previous_epoch_reserve"`
	CurrentEpochReserve    string `json:"current_epoch_reserve"`
}
