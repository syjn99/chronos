package reserves

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	mock "github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/rpc/testutil"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/network"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
)

func TestGetReserves_BadRequest(t *testing.T) {
	st, err := util.NewBeaconState()
	require.NoError(t, err)
	currentSlot := primitives.Slot(5000)
	require.NoError(t, st.SetSlot(currentSlot))
	mockChainService := &mock.ChainService{Optimistic: true}

	testCases := []struct {
		name         string
		path         string
		urlParams    map[string]string
		state        state.BeaconState
		errorMessage string
	}{
		{
			name:         "no state_id",
			path:         "/over/v1/beacon/states/{state_id}/reserves",
			urlParams:    map[string]string{},
			state:        nil,
			errorMessage: "state_id is required in URL params",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			s := &Server{
				FinalizationFetcher:   mockChainService,
				OptimisticModeFetcher: mockChainService,
				Stater:                &testutil.MockStater{BeaconState: testCase.state},
			}
			request := httptest.NewRequest("GET", testCase.path, nil)
			request = mux.SetURLVars(request, testCase.urlParams)
			writer := httptest.NewRecorder()
			writer.Body = &bytes.Buffer{}

			s.GetReserves(writer, request)
			assert.Equal(t, http.StatusBadRequest, writer.Code)
			e := &network.DefaultErrorJson{}
			require.NoError(t, json.Unmarshal(writer.Body.Bytes(), e))
			assert.Equal(t, http.StatusBadRequest, e.Code)
			assert.StringContains(t, testCase.errorMessage, e.Message)
		})
	}
}

func TestGetReserves(t *testing.T) {
	st, err := util.NewBeaconState()
	require.NoError(t, err)
	currentSlot := primitives.Slot(5000)
	require.NoError(t, st.SetSlot(currentSlot))
	mockChainService := &mock.ChainService{Optimistic: true}

	t.Run("get correct reserves data", func(t *testing.T) {
		params.SetupTestConfigCleanup(t)

		var (
			wantRewardAdjustmentFactor = uint64(200000)
			wantPreviousEpochReserve   = uint64(10000000)
			wantCurrentEpochReserve    = uint64(9000000)
		)

		require.NoError(t, st.SetRewardAdjustmentFactor(wantRewardAdjustmentFactor))
		require.NoError(t, st.SetPreviousEpochReserve(wantPreviousEpochReserve))
		require.NoError(t, st.SetCurrentEpochReserve(wantCurrentEpochReserve))

		s := &Server{
			FinalizationFetcher:   mockChainService,
			OptimisticModeFetcher: mockChainService,
			Stater:                &testutil.MockStater{BeaconState: st},
		}
		request := httptest.NewRequest(
			"GET", "/over/v1/beacon/states/{state_id}/reserves", nil)
		request = mux.SetURLVars(request, map[string]string{"state_id": "head"})
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}

		s.GetReserves(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		resp := &GetReservesResponse{}
		require.NoError(t, json.Unmarshal(writer.Body.Bytes(), resp))
		assert.Equal(t, true, resp.ExecutionOptimistic)
		assert.Equal(t, false, resp.Finalized)
		expectedReserves := &Reserves{
			RewardAdjustmentFactor: strconv.FormatUint(wantRewardAdjustmentFactor, 10),
			PreviousEpochReserve:   strconv.FormatUint(wantPreviousEpochReserve, 10),
			CurrentEpochReserve:    strconv.FormatUint(wantCurrentEpochReserve, 10),
		}
		require.DeepEqual(t, expectedReserves, resp.Data)
	})
}
