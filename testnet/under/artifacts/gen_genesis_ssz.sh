#!/bin/bash
BASE_DIR=$(pwd)
echo $BASE_DIR
KAIROS_DIR=$BASE_DIR/../../../../kairos
GETH_PATH=$KAIROS_DIR/build/bin/geth
$GETH_PATH --datadir $KAIROS_DIR/testnet/under/node-0 dump 0 &> $BASE_DIR/geth_output.txt
ETH1_GENESISHASH=$(grep "block=0 hash=" $BASE_DIR/geth_output.txt | awk -F'hash=' '{split($2, a, " "); print a[1]}' | sed 's/0x//')
rm $BASE_DIR/geth_output.txt

$BASE_DIR/eth2-256net-genesis phase0 \
    --config $BASE_DIR/config.yml \
    --state-output $BASE_DIR/genesis.ssz \
    --tranches-dir $BASE_DIR/tranches \
    --timestamp "1690458000" \
    --mnemonics $BASE_DIR/mnemonics.yaml \
    --eth1-block $ETH1_GENESISHASH

chmod 0600 $BASE_DIR/genesis.ssz