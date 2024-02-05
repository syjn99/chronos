package initialsync

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	initialSyncStartSlot = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "initial_sync_start_slot",
		Help: "The slot at which initial sync started.",
	})

	initialSyncReceivedBlocks = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "initial_sync_received_blocks",
		Help: "The amounts of blocks received during initial sync.",
	})

	initialSyncRunning = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "initial_sync_running",
		Help: "Indicates whether the node is currently waiting for minimum peers to start initial sync.",
	})
)

// SetInitialSyncStartSlotMetric sets the initial sync start slot
func setInitialSyncStartSlotMetric(slot float64) {
	initialSyncStartSlot.Set(slot)
}

// SetInitialSyncReceivedBlocksMetric sets the amount of blocks received during initial sync.
func setInitialSyncReceivedBlocksMetric(numBlocks float64) {
	initialSyncReceivedBlocks.Set(numBlocks)
}

// UpdateWaitForMinimumPeersMetric updates indicator whether we are waiting for minimum peers to start initial sync.
// 0 indicates that we are waiting for minimum peers to start initial sync.
// 1 indicates that we are not waiting for minimum peers to start initial sync.
func updateWaitForMinimumPeersMetric(isStart bool) {
	if isStart {
		initialSyncRunning.Set(0)
	} else {
		initialSyncRunning.Set(1)
	}
}

// UpdateInitialSyncReceivedBlocksMetric updates the amount of blocks received during initial sync.
func updateInitialSyncReceivedBlocksMetric(numBlocks float64) {
	initialSyncReceivedBlocks.Add(numBlocks)
}

// ResetInitialSyncMetrics resets all initial sync metrics.
func resetInitialSyncMetrics() {
	initialSyncStartSlot.Set(0)
	initialSyncReceivedBlocks.Set(0)
	initialSyncRunning.Set(0)
}
