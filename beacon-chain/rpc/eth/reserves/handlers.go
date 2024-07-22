package reserves

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v4/network"
)

// GetReserves get reserves data from the requested state.
// e.g. RewardAdjustmentFactor, PreviousEpochReserve, CurrentEpochReserve, etc.
func (s *Server) GetReserves(w http.ResponseWriter, r *http.Request) {
	// Retrieve beacon state
	stateId := mux.Vars(r)["state_id"]
	if stateId == "" {
		network.WriteError(w, &network.DefaultErrorJson{
			Message: "state_id is required in URL params",
			Code:    http.StatusBadRequest,
		})
		return
	}
	st, err := s.Stater.State(r.Context(), []byte(stateId))
	if err != nil {
		network.WriteError(w, handleWrapError(err, "could not retrieve state", http.StatusNotFound))
		return
	}
	// Get metadata for response
	isOptimistic, err := s.OptimisticModeFetcher.IsOptimistic(r.Context())
	if err != nil {
		network.WriteError(w, handleWrapError(err, "could not get optimistic mode info", http.StatusInternalServerError))
		return
	}
	root, err := helpers.BlockRootAtSlot(st, st.Slot()-1)
	if err != nil {
		network.WriteError(w, handleWrapError(err, "could not get block root", http.StatusInternalServerError))
		return
	}
	var blockRoot = [32]byte(root)
	isFinalized := s.FinalizationFetcher.IsFinalized(r.Context(), blockRoot)

	rewardAdjustmentFactor := st.RewardAdjustmentFactor()
	previousEpochReserve := st.PreviousEpochReserve()
	currentEpochReserve := st.CurrentEpochReserve()

	network.WriteJson(w, &GetReservesResponse{
		ExecutionOptimistic: isOptimistic,
		Finalized:           isFinalized,
		Data: &Reserves{
			RewardAdjustmentFactor: strconv.FormatUint(rewardAdjustmentFactor, 10),
			PreviousEpochReserve:   strconv.FormatUint(previousEpochReserve, 10),
			CurrentEpochReserve:    strconv.FormatUint(currentEpochReserve, 10),
		},
	})
}

func handleWrapError(err error, message string, code int) *network.DefaultErrorJson {
	return &network.DefaultErrorJson{
		Message: errors.Wrapf(err, message).Error(),
		Code:    code,
	}
}
