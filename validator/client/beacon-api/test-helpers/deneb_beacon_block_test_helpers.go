package test_helpers

import (
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
)

func GenerateProtoDenebBeaconBlockContents() *ethpb.BeaconBlockContentsDeneb {
	return &ethpb.BeaconBlockContentsDeneb{
		Block: &ethpb.BeaconBlockDeneb{
			Slot:          1,
			ProposerIndex: 2,
			ParentRoot:    FillByteSlice(32, 3),
			StateRoot:     FillByteSlice(32, 4),
			Body: &ethpb.BeaconBlockBodyDeneb{
				RandaoReveal: FillByteSlice(96, 5),
				Eth1Data: &ethpb.Eth1Data{
					DepositRoot:  FillByteSlice(32, 6),
					DepositCount: 7,
					BlockHash:    FillByteSlice(32, 8),
				},
				Graffiti: FillByteSlice(32, 9),
				ProposerSlashings: []*ethpb.ProposerSlashing{
					{
						Header_1: &ethpb.SignedBeaconBlockHeader{
							Header: &ethpb.BeaconBlockHeader{
								Slot:          10,
								ProposerIndex: 11,
								ParentRoot:    FillByteSlice(32, 12),
								StateRoot:     FillByteSlice(32, 13),
								BodyRoot:      FillByteSlice(32, 14),
							},
							Signature: FillByteSlice(96, 15),
						},
						Header_2: &ethpb.SignedBeaconBlockHeader{
							Header: &ethpb.BeaconBlockHeader{
								Slot:          16,
								ProposerIndex: 17,
								ParentRoot:    FillByteSlice(32, 18),
								StateRoot:     FillByteSlice(32, 19),
								BodyRoot:      FillByteSlice(32, 20),
							},
							Signature: FillByteSlice(96, 21),
						},
					},
					{
						Header_1: &ethpb.SignedBeaconBlockHeader{
							Header: &ethpb.BeaconBlockHeader{
								Slot:          22,
								ProposerIndex: 23,
								ParentRoot:    FillByteSlice(32, 24),
								StateRoot:     FillByteSlice(32, 25),
								BodyRoot:      FillByteSlice(32, 26),
							},
							Signature: FillByteSlice(96, 27),
						},
						Header_2: &ethpb.SignedBeaconBlockHeader{
							Header: &ethpb.BeaconBlockHeader{
								Slot:          28,
								ProposerIndex: 29,
								ParentRoot:    FillByteSlice(32, 30),
								StateRoot:     FillByteSlice(32, 31),
								BodyRoot:      FillByteSlice(32, 32),
							},
							Signature: FillByteSlice(96, 33),
						},
					},
				},
				AttesterSlashings: []*ethpb.AttesterSlashing{
					{
						Attestation_1: &ethpb.IndexedAttestation{
							AttestingIndices: []uint64{34, 35},
							Data: &ethpb.AttestationData{
								Slot:            36,
								CommitteeIndex:  37,
								BeaconBlockRoot: FillByteSlice(32, 38),
								Source: &ethpb.Checkpoint{
									Epoch: 39,
									Root:  FillByteSlice(32, 40),
								},
								Target: &ethpb.Checkpoint{
									Epoch: 41,
									Root:  FillByteSlice(32, 42),
								},
							},
							Signature: FillByteSlice(96, 43),
						},
						Attestation_2: &ethpb.IndexedAttestation{
							AttestingIndices: []uint64{44, 45},
							Data: &ethpb.AttestationData{
								Slot:            46,
								CommitteeIndex:  47,
								BeaconBlockRoot: FillByteSlice(32, 38),
								Source: &ethpb.Checkpoint{
									Epoch: 49,
									Root:  FillByteSlice(32, 50),
								},
								Target: &ethpb.Checkpoint{
									Epoch: 51,
									Root:  FillByteSlice(32, 52),
								},
							},
							Signature: FillByteSlice(96, 53),
						},
					},
					{
						Attestation_1: &ethpb.IndexedAttestation{
							AttestingIndices: []uint64{54, 55},
							Data: &ethpb.AttestationData{
								Slot:            56,
								CommitteeIndex:  57,
								BeaconBlockRoot: FillByteSlice(32, 38),
								Source: &ethpb.Checkpoint{
									Epoch: 59,
									Root:  FillByteSlice(32, 60),
								},
								Target: &ethpb.Checkpoint{
									Epoch: 61,
									Root:  FillByteSlice(32, 62),
								},
							},
							Signature: FillByteSlice(96, 63),
						},
						Attestation_2: &ethpb.IndexedAttestation{
							AttestingIndices: []uint64{64, 65},
							Data: &ethpb.AttestationData{
								Slot:            66,
								CommitteeIndex:  67,
								BeaconBlockRoot: FillByteSlice(32, 38),
								Source: &ethpb.Checkpoint{
									Epoch: 69,
									Root:  FillByteSlice(32, 70),
								},
								Target: &ethpb.Checkpoint{
									Epoch: 71,
									Root:  FillByteSlice(32, 72),
								},
							},
							Signature: FillByteSlice(96, 73),
						},
					},
				},
				Attestations: []*ethpb.Attestation{
					{
						AggregationBits: FillByteSlice(4, 74),
						Data: &ethpb.AttestationData{
							Slot:            75,
							CommitteeIndex:  76,
							BeaconBlockRoot: FillByteSlice(32, 38),
							Source: &ethpb.Checkpoint{
								Epoch: 78,
								Root:  FillByteSlice(32, 79),
							},
							Target: &ethpb.Checkpoint{
								Epoch: 80,
								Root:  FillByteSlice(32, 81),
							},
						},
						Signature: FillByteSlice(96, 82),
					},
					{
						AggregationBits: FillByteSlice(4, 83),
						Data: &ethpb.AttestationData{
							Slot:            84,
							CommitteeIndex:  85,
							BeaconBlockRoot: FillByteSlice(32, 38),
							Source: &ethpb.Checkpoint{
								Epoch: 87,
								Root:  FillByteSlice(32, 88),
							},
							Target: &ethpb.Checkpoint{
								Epoch: 89,
								Root:  FillByteSlice(32, 90),
							},
						},
						Signature: FillByteSlice(96, 91),
					},
				},
				Deposits: []*ethpb.Deposit{
					{
						Proof: FillByteArraySlice(33, FillByteSlice(32, 92)),
						Data: &ethpb.Deposit_Data{
							PublicKey:             FillByteSlice(48, 94),
							WithdrawalCredentials: FillByteSlice(32, 95),
							Amount:                96,
							Signature:             FillByteSlice(96, 97),
						},
					},
					{
						Proof: FillByteArraySlice(33, FillByteSlice(32, 98)),
						Data: &ethpb.Deposit_Data{
							PublicKey:             FillByteSlice(48, 100),
							WithdrawalCredentials: FillByteSlice(32, 101),
							Amount:                102,
							Signature:             FillByteSlice(96, 103),
						},
					},
				},
				VoluntaryExits: []*ethpb.SignedVoluntaryExit{
					{
						Exit: &ethpb.VoluntaryExit{
							Epoch:          104,
							ValidatorIndex: 105,
						},
						Signature: FillByteSlice(96, 106),
					},
					{
						Exit: &ethpb.VoluntaryExit{
							Epoch:          107,
							ValidatorIndex: 108,
						},
						Signature: FillByteSlice(96, 109),
					},
				},
				SyncAggregate: &ethpb.SyncAggregate{
					SyncCommitteeBits:      FillByteSlice(64, 110),
					SyncCommitteeSignature: FillByteSlice(96, 111),
				},
				BailOuts: []*ethpb.BailOut{
					{
						ValidatorIndex: 128,
					},
					{
						ValidatorIndex: 129,
					},
				},
				ExecutionPayload: &enginev1.ExecutionPayloadDeneb{
					ParentHash:     FillByteSlice(32, 112),
					FeeRecipient:   FillByteSlice(20, 113),
					StateRoot:      FillByteSlice(32, 114),
					CheckpointRoot: FillByteSlice(32, 115),
					ReceiptsRoot:   FillByteSlice(32, 116),
					LogsBloom:      FillByteSlice(256, 117),
					PrevRandao:     FillByteSlice(32, 118),
					BlockNumber:    119,
					GasLimit:       120,
					GasUsed:        121,
					Timestamp:      122,
					ExtraData:      FillByteSlice(32, 123),
					BaseFeePerGas:  FillByteSlice(32, 124),
					BlockHash:      FillByteSlice(32, 125),
					Transactions: [][]byte{
						FillByteSlice(32, 126),
						FillByteSlice(32, 127),
					},
					Withdrawals: []*enginev1.Withdrawal{
						{
							Index:          128,
							ValidatorIndex: 129,
							Address:        FillByteSlice(20, 130),
							Amount:         131,
						},
						{
							Index:          132,
							ValidatorIndex: 133,
							Address:        FillByteSlice(20, 134),
							Amount:         135,
						},
					},
					BlobGasUsed:   136,
					ExcessBlobGas: 137,
				},
				BlsToExecutionChanges: []*ethpb.SignedBLSToExecutionChange{
					{
						Message: &ethpb.BLSToExecutionChange{
							ValidatorIndex:     138,
							FromBlsPubkey:      FillByteSlice(48, 139),
							ToExecutionAddress: FillByteSlice(20, 140),
						},
						Signature: FillByteSlice(96, 141),
					},
					{
						Message: &ethpb.BLSToExecutionChange{
							ValidatorIndex:     142,
							FromBlsPubkey:      FillByteSlice(48, 143),
							ToExecutionAddress: FillByteSlice(20, 144),
						},
						Signature: FillByteSlice(96, 145),
					},
				},
				BlobKzgCommitments: [][]byte{FillByteSlice(48, 146), FillByteSlice(48, 147)},
			},
		},
		KzgProofs: [][]byte{FillByteSlice(48, 148)},
		Blobs:     [][]byte{FillByteSlice(131072, 149)},
	}
}

func GenerateProtoBlindedDenebBeaconBlock() *ethpb.BlindedBeaconBlockDeneb {
	return &ethpb.BlindedBeaconBlockDeneb{
		Slot:          1,
		ProposerIndex: 2,
		ParentRoot:    FillByteSlice(32, 3),
		StateRoot:     FillByteSlice(32, 4),
		Body: &ethpb.BlindedBeaconBlockBodyDeneb{
			RandaoReveal: FillByteSlice(96, 5),
			Eth1Data: &ethpb.Eth1Data{
				DepositRoot:  FillByteSlice(32, 6),
				DepositCount: 7,
				BlockHash:    FillByteSlice(32, 8),
			},
			Graffiti: FillByteSlice(32, 9),
			ProposerSlashings: []*ethpb.ProposerSlashing{
				{
					Header_1: &ethpb.SignedBeaconBlockHeader{
						Header: &ethpb.BeaconBlockHeader{
							Slot:          10,
							ProposerIndex: 11,
							ParentRoot:    FillByteSlice(32, 12),
							StateRoot:     FillByteSlice(32, 13),
							BodyRoot:      FillByteSlice(32, 14),
						},
						Signature: FillByteSlice(96, 15),
					},
					Header_2: &ethpb.SignedBeaconBlockHeader{
						Header: &ethpb.BeaconBlockHeader{
							Slot:          16,
							ProposerIndex: 17,
							ParentRoot:    FillByteSlice(32, 18),
							StateRoot:     FillByteSlice(32, 19),
							BodyRoot:      FillByteSlice(32, 20),
						},
						Signature: FillByteSlice(96, 21),
					},
				},
				{
					Header_1: &ethpb.SignedBeaconBlockHeader{
						Header: &ethpb.BeaconBlockHeader{
							Slot:          22,
							ProposerIndex: 23,
							ParentRoot:    FillByteSlice(32, 24),
							StateRoot:     FillByteSlice(32, 25),
							BodyRoot:      FillByteSlice(32, 26),
						},
						Signature: FillByteSlice(96, 27),
					},
					Header_2: &ethpb.SignedBeaconBlockHeader{
						Header: &ethpb.BeaconBlockHeader{
							Slot:          28,
							ProposerIndex: 29,
							ParentRoot:    FillByteSlice(32, 30),
							StateRoot:     FillByteSlice(32, 31),
							BodyRoot:      FillByteSlice(32, 32),
						},
						Signature: FillByteSlice(96, 33),
					},
				},
			},
			AttesterSlashings: []*ethpb.AttesterSlashing{
				{
					Attestation_1: &ethpb.IndexedAttestation{
						AttestingIndices: []uint64{34, 35},
						Data: &ethpb.AttestationData{
							Slot:            36,
							CommitteeIndex:  37,
							BeaconBlockRoot: FillByteSlice(32, 38),
							Source: &ethpb.Checkpoint{
								Epoch: 39,
								Root:  FillByteSlice(32, 40),
							},
							Target: &ethpb.Checkpoint{
								Epoch: 41,
								Root:  FillByteSlice(32, 42),
							},
						},
						Signature: FillByteSlice(96, 43),
					},
					Attestation_2: &ethpb.IndexedAttestation{
						AttestingIndices: []uint64{44, 45},
						Data: &ethpb.AttestationData{
							Slot:            46,
							CommitteeIndex:  47,
							BeaconBlockRoot: FillByteSlice(32, 38),
							Source: &ethpb.Checkpoint{
								Epoch: 49,
								Root:  FillByteSlice(32, 50),
							},
							Target: &ethpb.Checkpoint{
								Epoch: 51,
								Root:  FillByteSlice(32, 52),
							},
						},
						Signature: FillByteSlice(96, 53),
					},
				},
				{
					Attestation_1: &ethpb.IndexedAttestation{
						AttestingIndices: []uint64{54, 55},
						Data: &ethpb.AttestationData{
							Slot:            56,
							CommitteeIndex:  57,
							BeaconBlockRoot: FillByteSlice(32, 38),
							Source: &ethpb.Checkpoint{
								Epoch: 59,
								Root:  FillByteSlice(32, 60),
							},
							Target: &ethpb.Checkpoint{
								Epoch: 61,
								Root:  FillByteSlice(32, 62),
							},
						},
						Signature: FillByteSlice(96, 63),
					},
					Attestation_2: &ethpb.IndexedAttestation{
						AttestingIndices: []uint64{64, 65},
						Data: &ethpb.AttestationData{
							Slot:            66,
							CommitteeIndex:  67,
							BeaconBlockRoot: FillByteSlice(32, 38),
							Source: &ethpb.Checkpoint{
								Epoch: 69,
								Root:  FillByteSlice(32, 70),
							},
							Target: &ethpb.Checkpoint{
								Epoch: 71,
								Root:  FillByteSlice(32, 72),
							},
						},
						Signature: FillByteSlice(96, 73),
					},
				},
			},
			Attestations: []*ethpb.Attestation{
				{
					AggregationBits: FillByteSlice(4, 74),
					Data: &ethpb.AttestationData{
						Slot:            75,
						CommitteeIndex:  76,
						BeaconBlockRoot: FillByteSlice(32, 38),
						Source: &ethpb.Checkpoint{
							Epoch: 78,
							Root:  FillByteSlice(32, 79),
						},
						Target: &ethpb.Checkpoint{
							Epoch: 80,
							Root:  FillByteSlice(32, 81),
						},
					},
					Signature: FillByteSlice(96, 82),
				},
				{
					AggregationBits: FillByteSlice(4, 83),
					Data: &ethpb.AttestationData{
						Slot:            84,
						CommitteeIndex:  85,
						BeaconBlockRoot: FillByteSlice(32, 38),
						Source: &ethpb.Checkpoint{
							Epoch: 87,
							Root:  FillByteSlice(32, 88),
						},
						Target: &ethpb.Checkpoint{
							Epoch: 89,
							Root:  FillByteSlice(32, 90),
						},
					},
					Signature: FillByteSlice(96, 91),
				},
			},
			Deposits: []*ethpb.Deposit{
				{
					Proof: FillByteArraySlice(33, FillByteSlice(32, 92)),
					Data: &ethpb.Deposit_Data{
						PublicKey:             FillByteSlice(48, 94),
						WithdrawalCredentials: FillByteSlice(32, 95),
						Amount:                96,
						Signature:             FillByteSlice(96, 97),
					},
				},
				{
					Proof: FillByteArraySlice(33, FillByteSlice(32, 98)),
					Data: &ethpb.Deposit_Data{
						PublicKey:             FillByteSlice(48, 100),
						WithdrawalCredentials: FillByteSlice(32, 101),
						Amount:                102,
						Signature:             FillByteSlice(96, 103),
					},
				},
			},
			VoluntaryExits: []*ethpb.SignedVoluntaryExit{
				{
					Exit: &ethpb.VoluntaryExit{
						Epoch:          104,
						ValidatorIndex: 105,
					},
					Signature: FillByteSlice(96, 106),
				},
				{
					Exit: &ethpb.VoluntaryExit{
						Epoch:          107,
						ValidatorIndex: 108,
					},
					Signature: FillByteSlice(96, 109),
				},
			},
			SyncAggregate: &ethpb.SyncAggregate{
				SyncCommitteeBits:      FillByteSlice(64, 110),
				SyncCommitteeSignature: FillByteSlice(96, 111),
			},
			BailOuts: []*ethpb.BailOut{
				{
					ValidatorIndex: 128,
				},
				{
					ValidatorIndex: 129,
				},
			},
			ExecutionPayloadHeader: &enginev1.ExecutionPayloadHeaderDeneb{
				ParentHash:       FillByteSlice(32, 112),
				FeeRecipient:     FillByteSlice(20, 113),
				StateRoot:        FillByteSlice(32, 114),
				CheckpointRoot:   FillByteSlice(32, 115),
				ReceiptsRoot:     FillByteSlice(32, 116),
				LogsBloom:        FillByteSlice(256, 117),
				PrevRandao:       FillByteSlice(32, 118),
				BlockNumber:      119,
				GasLimit:         120,
				GasUsed:          121,
				Timestamp:        122,
				ExtraData:        FillByteSlice(32, 123),
				BaseFeePerGas:    FillByteSlice(32, 124),
				BlockHash:        FillByteSlice(32, 125),
				TransactionsRoot: FillByteSlice(32, 126),
				WithdrawalsRoot:  FillByteSlice(32, 127),
				BlobGasUsed:      128,
				ExcessBlobGas:    129,
			},
			BlsToExecutionChanges: []*ethpb.SignedBLSToExecutionChange{
				{
					Message: &ethpb.BLSToExecutionChange{
						ValidatorIndex:     130,
						FromBlsPubkey:      FillByteSlice(48, 131),
						ToExecutionAddress: FillByteSlice(20, 132),
					},
					Signature: FillByteSlice(96, 133),
				},
				{
					Message: &ethpb.BLSToExecutionChange{
						ValidatorIndex:     134,
						FromBlsPubkey:      FillByteSlice(48, 135),
						ToExecutionAddress: FillByteSlice(20, 136),
					},
					Signature: FillByteSlice(96, 137),
				},
			},
			BlobKzgCommitments: [][]byte{FillByteSlice(48, 138), FillByteSlice(48, 139)},
		},
	}
}

func GenerateJsonDenebBeaconBlockContents() *structs.BeaconBlockContentsDeneb {
	return &structs.BeaconBlockContentsDeneb{
		Block: &structs.BeaconBlockDeneb{
			Slot:          "1",
			ProposerIndex: "2",
			ParentRoot:    FillEncodedByteSlice(32, 3),
			StateRoot:     FillEncodedByteSlice(32, 4),
			Body: &structs.BeaconBlockBodyDeneb{
				RandaoReveal: FillEncodedByteSlice(96, 5),
				Eth1Data: &structs.Eth1Data{
					DepositRoot:  FillEncodedByteSlice(32, 6),
					DepositCount: "7",
					BlockHash:    FillEncodedByteSlice(32, 8),
				},
				Graffiti: FillEncodedByteSlice(32, 9),
				ProposerSlashings: []*structs.ProposerSlashing{
					{
						SignedHeader1: &structs.SignedBeaconBlockHeader{
							Message: &structs.BeaconBlockHeader{
								Slot:          "10",
								ProposerIndex: "11",
								ParentRoot:    FillEncodedByteSlice(32, 12),
								StateRoot:     FillEncodedByteSlice(32, 13),
								BodyRoot:      FillEncodedByteSlice(32, 14),
							},
							Signature: FillEncodedByteSlice(96, 15),
						},
						SignedHeader2: &structs.SignedBeaconBlockHeader{
							Message: &structs.BeaconBlockHeader{
								Slot:          "16",
								ProposerIndex: "17",
								ParentRoot:    FillEncodedByteSlice(32, 18),
								StateRoot:     FillEncodedByteSlice(32, 19),
								BodyRoot:      FillEncodedByteSlice(32, 20),
							},
							Signature: FillEncodedByteSlice(96, 21),
						},
					},
					{
						SignedHeader1: &structs.SignedBeaconBlockHeader{
							Message: &structs.BeaconBlockHeader{
								Slot:          "22",
								ProposerIndex: "23",
								ParentRoot:    FillEncodedByteSlice(32, 24),
								StateRoot:     FillEncodedByteSlice(32, 25),
								BodyRoot:      FillEncodedByteSlice(32, 26),
							},
							Signature: FillEncodedByteSlice(96, 27),
						},
						SignedHeader2: &structs.SignedBeaconBlockHeader{
							Message: &structs.BeaconBlockHeader{
								Slot:          "28",
								ProposerIndex: "29",
								ParentRoot:    FillEncodedByteSlice(32, 30),
								StateRoot:     FillEncodedByteSlice(32, 31),
								BodyRoot:      FillEncodedByteSlice(32, 32),
							},
							Signature: FillEncodedByteSlice(96, 33),
						},
					},
				},
				AttesterSlashings: []*structs.AttesterSlashing{
					{
						Attestation1: &structs.IndexedAttestation{
							AttestingIndices: []string{"34", "35"},
							Data: &structs.AttestationData{
								Slot:            "36",
								CommitteeIndex:  "37",
								BeaconBlockRoot: FillEncodedByteSlice(32, 38),
								Source: &structs.Checkpoint{
									Epoch: "39",
									Root:  FillEncodedByteSlice(32, 40),
								},
								Target: &structs.Checkpoint{
									Epoch: "41",
									Root:  FillEncodedByteSlice(32, 42),
								},
							},
							Signature: FillEncodedByteSlice(96, 43),
						},
						Attestation2: &structs.IndexedAttestation{
							AttestingIndices: []string{"44", "45"},
							Data: &structs.AttestationData{
								Slot:            "46",
								CommitteeIndex:  "47",
								BeaconBlockRoot: FillEncodedByteSlice(32, 38),
								Source: &structs.Checkpoint{
									Epoch: "49",
									Root:  FillEncodedByteSlice(32, 50),
								},
								Target: &structs.Checkpoint{
									Epoch: "51",
									Root:  FillEncodedByteSlice(32, 52),
								},
							},
							Signature: FillEncodedByteSlice(96, 53),
						},
					},
					{
						Attestation1: &structs.IndexedAttestation{
							AttestingIndices: []string{"54", "55"},
							Data: &structs.AttestationData{
								Slot:            "56",
								CommitteeIndex:  "57",
								BeaconBlockRoot: FillEncodedByteSlice(32, 38),
								Source: &structs.Checkpoint{
									Epoch: "59",
									Root:  FillEncodedByteSlice(32, 60),
								},
								Target: &structs.Checkpoint{
									Epoch: "61",
									Root:  FillEncodedByteSlice(32, 62),
								},
							},
							Signature: FillEncodedByteSlice(96, 63),
						},
						Attestation2: &structs.IndexedAttestation{
							AttestingIndices: []string{"64", "65"},
							Data: &structs.AttestationData{
								Slot:            "66",
								CommitteeIndex:  "67",
								BeaconBlockRoot: FillEncodedByteSlice(32, 38),
								Source: &structs.Checkpoint{
									Epoch: "69",
									Root:  FillEncodedByteSlice(32, 70),
								},
								Target: &structs.Checkpoint{
									Epoch: "71",
									Root:  FillEncodedByteSlice(32, 72),
								},
							},
							Signature: FillEncodedByteSlice(96, 73),
						},
					},
				},
				Attestations: []*structs.Attestation{
					{
						AggregationBits: FillEncodedByteSlice(4, 74),
						Data: &structs.AttestationData{
							Slot:            "75",
							CommitteeIndex:  "76",
							BeaconBlockRoot: FillEncodedByteSlice(32, 38),
							Source: &structs.Checkpoint{
								Epoch: "78",
								Root:  FillEncodedByteSlice(32, 79),
							},
							Target: &structs.Checkpoint{
								Epoch: "80",
								Root:  FillEncodedByteSlice(32, 81),
							},
						},
						Signature: FillEncodedByteSlice(96, 82),
					},
					{
						AggregationBits: FillEncodedByteSlice(4, 83),
						Data: &structs.AttestationData{
							Slot:            "84",
							CommitteeIndex:  "85",
							BeaconBlockRoot: FillEncodedByteSlice(32, 38),
							Source: &structs.Checkpoint{
								Epoch: "87",
								Root:  FillEncodedByteSlice(32, 88),
							},
							Target: &structs.Checkpoint{
								Epoch: "89",
								Root:  FillEncodedByteSlice(32, 90),
							},
						},
						Signature: FillEncodedByteSlice(96, 91),
					},
				},
				Deposits: []*structs.Deposit{
					{
						Proof: FillEncodedByteArraySlice(33, FillEncodedByteSlice(32, 92)),
						Data: &structs.DepositData{
							Pubkey:                FillEncodedByteSlice(48, 94),
							WithdrawalCredentials: FillEncodedByteSlice(32, 95),
							Amount:                "96",
							Signature:             FillEncodedByteSlice(96, 97),
						},
					},
					{
						Proof: FillEncodedByteArraySlice(33, FillEncodedByteSlice(32, 98)),
						Data: &structs.DepositData{
							Pubkey:                FillEncodedByteSlice(48, 100),
							WithdrawalCredentials: FillEncodedByteSlice(32, 101),
							Amount:                "102",
							Signature:             FillEncodedByteSlice(96, 103),
						},
					},
				},
				VoluntaryExits: []*structs.SignedVoluntaryExit{
					{
						Message: &structs.VoluntaryExit{
							Epoch:          "104",
							ValidatorIndex: "105",
						},
						Signature: FillEncodedByteSlice(96, 106),
					},
					{
						Message: &structs.VoluntaryExit{
							Epoch:          "107",
							ValidatorIndex: "108",
						},
						Signature: FillEncodedByteSlice(96, 109),
					},
				},
				SyncAggregate: &structs.SyncAggregate{
					SyncCommitteeBits:      FillEncodedByteSlice(64, 110),
					SyncCommitteeSignature: FillEncodedByteSlice(96, 111),
				},
				BailOuts: []*structs.BailOut{
					{
						ValidatorIndex: "128",
					},
					{
						ValidatorIndex: "129",
					},
				},
				ExecutionPayload: &structs.ExecutionPayloadDeneb{
					ParentHash:     FillEncodedByteSlice(32, 112),
					FeeRecipient:   FillEncodedByteSlice(20, 113),
					StateRoot:      FillEncodedByteSlice(32, 114),
					CheckpointRoot: FillEncodedByteSlice(32, 115),
					ReceiptsRoot:   FillEncodedByteSlice(32, 116),
					LogsBloom:      FillEncodedByteSlice(256, 117),
					PrevRandao:     FillEncodedByteSlice(32, 118),
					BlockNumber:    "119",
					GasLimit:       "120",
					GasUsed:        "121",
					Timestamp:      "122",
					ExtraData:      FillEncodedByteSlice(32, 123),
					BaseFeePerGas:  bytesutil.LittleEndianBytesToBigInt(FillByteSlice(32, 124)).String(),
					BlockHash:      FillEncodedByteSlice(32, 125),
					Transactions: []string{
						FillEncodedByteSlice(32, 126),
						FillEncodedByteSlice(32, 127),
					},
					Withdrawals: []*structs.Withdrawal{
						{
							WithdrawalIndex:  "128",
							ValidatorIndex:   "129",
							ExecutionAddress: FillEncodedByteSlice(20, 130),
							Amount:           "131",
						},
						{
							WithdrawalIndex:  "132",
							ValidatorIndex:   "133",
							ExecutionAddress: FillEncodedByteSlice(20, 134),
							Amount:           "135",
						},
					},
					BlobGasUsed:   "136",
					ExcessBlobGas: "137",
				},
				BLSToExecutionChanges: []*structs.SignedBLSToExecutionChange{
					{
						Message: &structs.BLSToExecutionChange{
							ValidatorIndex:     "138",
							FromBLSPubkey:      FillEncodedByteSlice(48, 139),
							ToExecutionAddress: FillEncodedByteSlice(20, 140),
						},
						Signature: FillEncodedByteSlice(96, 141),
					},
					{
						Message: &structs.BLSToExecutionChange{
							ValidatorIndex:     "142",
							FromBLSPubkey:      FillEncodedByteSlice(48, 143),
							ToExecutionAddress: FillEncodedByteSlice(20, 144),
						},
						Signature: FillEncodedByteSlice(96, 145),
					},
				},
				BlobKzgCommitments: []string{FillEncodedByteSlice(48, 146), FillEncodedByteSlice(48, 147)},
			},
		},
		KzgProofs: []string{FillEncodedByteSlice(48, 148)},
		Blobs:     []string{FillEncodedByteSlice(131072, 149)},
	}
}

func GenerateJsonBlindedDenebBeaconBlock() *structs.BlindedBeaconBlockDeneb {
	return &structs.BlindedBeaconBlockDeneb{
		Slot:          "1",
		ProposerIndex: "2",
		ParentRoot:    FillEncodedByteSlice(32, 3),
		StateRoot:     FillEncodedByteSlice(32, 4),
		Body: &structs.BlindedBeaconBlockBodyDeneb{
			RandaoReveal: FillEncodedByteSlice(96, 5),
			Eth1Data: &structs.Eth1Data{
				DepositRoot:  FillEncodedByteSlice(32, 6),
				DepositCount: "7",
				BlockHash:    FillEncodedByteSlice(32, 8),
			},
			Graffiti: FillEncodedByteSlice(32, 9),
			ProposerSlashings: []*structs.ProposerSlashing{
				{
					SignedHeader1: &structs.SignedBeaconBlockHeader{
						Message: &structs.BeaconBlockHeader{
							Slot:          "10",
							ProposerIndex: "11",
							ParentRoot:    FillEncodedByteSlice(32, 12),
							StateRoot:     FillEncodedByteSlice(32, 13),
							BodyRoot:      FillEncodedByteSlice(32, 14),
						},
						Signature: FillEncodedByteSlice(96, 15),
					},
					SignedHeader2: &structs.SignedBeaconBlockHeader{
						Message: &structs.BeaconBlockHeader{
							Slot:          "16",
							ProposerIndex: "17",
							ParentRoot:    FillEncodedByteSlice(32, 18),
							StateRoot:     FillEncodedByteSlice(32, 19),
							BodyRoot:      FillEncodedByteSlice(32, 20),
						},
						Signature: FillEncodedByteSlice(96, 21),
					},
				},
				{
					SignedHeader1: &structs.SignedBeaconBlockHeader{
						Message: &structs.BeaconBlockHeader{
							Slot:          "22",
							ProposerIndex: "23",
							ParentRoot:    FillEncodedByteSlice(32, 24),
							StateRoot:     FillEncodedByteSlice(32, 25),
							BodyRoot:      FillEncodedByteSlice(32, 26),
						},
						Signature: FillEncodedByteSlice(96, 27),
					},
					SignedHeader2: &structs.SignedBeaconBlockHeader{
						Message: &structs.BeaconBlockHeader{
							Slot:          "28",
							ProposerIndex: "29",
							ParentRoot:    FillEncodedByteSlice(32, 30),
							StateRoot:     FillEncodedByteSlice(32, 31),
							BodyRoot:      FillEncodedByteSlice(32, 32),
						},
						Signature: FillEncodedByteSlice(96, 33),
					},
				},
			},
			AttesterSlashings: []*structs.AttesterSlashing{
				{
					Attestation1: &structs.IndexedAttestation{
						AttestingIndices: []string{"34", "35"},
						Data: &structs.AttestationData{
							Slot:            "36",
							CommitteeIndex:  "37",
							BeaconBlockRoot: FillEncodedByteSlice(32, 38),
							Source: &structs.Checkpoint{
								Epoch: "39",
								Root:  FillEncodedByteSlice(32, 40),
							},
							Target: &structs.Checkpoint{
								Epoch: "41",
								Root:  FillEncodedByteSlice(32, 42),
							},
						},
						Signature: FillEncodedByteSlice(96, 43),
					},
					Attestation2: &structs.IndexedAttestation{
						AttestingIndices: []string{"44", "45"},
						Data: &structs.AttestationData{
							Slot:            "46",
							CommitteeIndex:  "47",
							BeaconBlockRoot: FillEncodedByteSlice(32, 38),
							Source: &structs.Checkpoint{
								Epoch: "49",
								Root:  FillEncodedByteSlice(32, 50),
							},
							Target: &structs.Checkpoint{
								Epoch: "51",
								Root:  FillEncodedByteSlice(32, 52),
							},
						},
						Signature: FillEncodedByteSlice(96, 53),
					},
				},
				{
					Attestation1: &structs.IndexedAttestation{
						AttestingIndices: []string{"54", "55"},
						Data: &structs.AttestationData{
							Slot:            "56",
							CommitteeIndex:  "57",
							BeaconBlockRoot: FillEncodedByteSlice(32, 38),
							Source: &structs.Checkpoint{
								Epoch: "59",
								Root:  FillEncodedByteSlice(32, 60),
							},
							Target: &structs.Checkpoint{
								Epoch: "61",
								Root:  FillEncodedByteSlice(32, 62),
							},
						},
						Signature: FillEncodedByteSlice(96, 63),
					},
					Attestation2: &structs.IndexedAttestation{
						AttestingIndices: []string{"64", "65"},
						Data: &structs.AttestationData{
							Slot:            "66",
							CommitteeIndex:  "67",
							BeaconBlockRoot: FillEncodedByteSlice(32, 38),
							Source: &structs.Checkpoint{
								Epoch: "69",
								Root:  FillEncodedByteSlice(32, 70),
							},
							Target: &structs.Checkpoint{
								Epoch: "71",
								Root:  FillEncodedByteSlice(32, 72),
							},
						},
						Signature: FillEncodedByteSlice(96, 73),
					},
				},
			},
			Attestations: []*structs.Attestation{
				{
					AggregationBits: FillEncodedByteSlice(4, 74),
					Data: &structs.AttestationData{
						Slot:            "75",
						CommitteeIndex:  "76",
						BeaconBlockRoot: FillEncodedByteSlice(32, 38),
						Source: &structs.Checkpoint{
							Epoch: "78",
							Root:  FillEncodedByteSlice(32, 79),
						},
						Target: &structs.Checkpoint{
							Epoch: "80",
							Root:  FillEncodedByteSlice(32, 81),
						},
					},
					Signature: FillEncodedByteSlice(96, 82),
				},
				{
					AggregationBits: FillEncodedByteSlice(4, 83),
					Data: &structs.AttestationData{
						Slot:            "84",
						CommitteeIndex:  "85",
						BeaconBlockRoot: FillEncodedByteSlice(32, 38),
						Source: &structs.Checkpoint{
							Epoch: "87",
							Root:  FillEncodedByteSlice(32, 88),
						},
						Target: &structs.Checkpoint{
							Epoch: "89",
							Root:  FillEncodedByteSlice(32, 90),
						},
					},
					Signature: FillEncodedByteSlice(96, 91),
				},
			},
			Deposits: []*structs.Deposit{
				{
					Proof: FillEncodedByteArraySlice(33, FillEncodedByteSlice(32, 92)),
					Data: &structs.DepositData{
						Pubkey:                FillEncodedByteSlice(48, 94),
						WithdrawalCredentials: FillEncodedByteSlice(32, 95),
						Amount:                "96",
						Signature:             FillEncodedByteSlice(96, 97),
					},
				},
				{
					Proof: FillEncodedByteArraySlice(33, FillEncodedByteSlice(32, 98)),
					Data: &structs.DepositData{
						Pubkey:                FillEncodedByteSlice(48, 100),
						WithdrawalCredentials: FillEncodedByteSlice(32, 101),
						Amount:                "102",
						Signature:             FillEncodedByteSlice(96, 103),
					},
				},
			},
			VoluntaryExits: []*structs.SignedVoluntaryExit{
				{
					Message: &structs.VoluntaryExit{
						Epoch:          "104",
						ValidatorIndex: "105",
					},
					Signature: FillEncodedByteSlice(96, 106),
				},
				{
					Message: &structs.VoluntaryExit{
						Epoch:          "107",
						ValidatorIndex: "108",
					},
					Signature: FillEncodedByteSlice(96, 109),
				},
			},
			SyncAggregate: &structs.SyncAggregate{
				SyncCommitteeBits:      FillEncodedByteSlice(64, 110),
				SyncCommitteeSignature: FillEncodedByteSlice(96, 111),
			},
			BailOuts: []*structs.BailOut{
				{
					ValidatorIndex: "128",
				},
				{
					ValidatorIndex: "129",
				},
			},
			ExecutionPayloadHeader: &structs.ExecutionPayloadHeaderDeneb{
				ParentHash:       FillEncodedByteSlice(32, 112),
				FeeRecipient:     FillEncodedByteSlice(20, 113),
				StateRoot:        FillEncodedByteSlice(32, 114),
				CheckpointRoot:   FillEncodedByteSlice(32, 115),
				ReceiptsRoot:     FillEncodedByteSlice(32, 116),
				LogsBloom:        FillEncodedByteSlice(256, 117),
				PrevRandao:       FillEncodedByteSlice(32, 118),
				BlockNumber:      "119",
				GasLimit:         "120",
				GasUsed:          "121",
				Timestamp:        "122",
				ExtraData:        FillEncodedByteSlice(32, 123),
				BaseFeePerGas:    bytesutil.LittleEndianBytesToBigInt(FillByteSlice(32, 124)).String(),
				BlockHash:        FillEncodedByteSlice(32, 125),
				TransactionsRoot: FillEncodedByteSlice(32, 126),
				WithdrawalsRoot:  FillEncodedByteSlice(32, 127),
				BlobGasUsed:      "128",
				ExcessBlobGas:    "129",
			},
			BLSToExecutionChanges: []*structs.SignedBLSToExecutionChange{
				{
					Message: &structs.BLSToExecutionChange{
						ValidatorIndex:     "130",
						FromBLSPubkey:      FillEncodedByteSlice(48, 131),
						ToExecutionAddress: FillEncodedByteSlice(20, 132),
					},
					Signature: FillEncodedByteSlice(96, 133),
				},
				{
					Message: &structs.BLSToExecutionChange{
						ValidatorIndex:     "134",
						FromBLSPubkey:      FillEncodedByteSlice(48, 135),
						ToExecutionAddress: FillEncodedByteSlice(20, 136),
					},
					Signature: FillEncodedByteSlice(96, 137),
				},
			},
			BlobKzgCommitments: []string{FillEncodedByteSlice(48, 138), FillEncodedByteSlice(48, 139)},
		},
	}
}
