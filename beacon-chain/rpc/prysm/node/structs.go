package node

import "time"

type AddrRequest struct {
	Addr string `json:"addr"`
}

type PeersResponse struct {
	Peers []*Peer `json:"Peers"`
}

type Peer struct {
	PeerID             string `json:"peer_id"`
	Enr                string `json:"enr"`
	LastSeenP2PAddress string `json:"last_seen_p2p_address"`
	State              string `json:"state"`
	Direction          string `json:"direction"`
}

type PeerDetailInfoResponse struct {
	PeerID       string `json:"peer_id"`
	Enr          string `json:"enr"`
	Address      string `json:"address"`
	IpTrackerCnt uint64 `json:"ip_tracker_cnt"`

	// Scorers internal Data
	BadResponses         int       `json:"bad_responses"`
	ProcessedBlocks      uint64    `json:"processed_blocks"`
	BlockProviderUpdated time.Time `json:"block_provider_updated"`
	// Gossip Scoring data.
	GossipScore      string `json:"gossip_score"`
	BehaviourPenalty string `json:"behaviour_penalty"`
}

type EpochReward struct {
	Reward uint64 `json:"reward"`
}
