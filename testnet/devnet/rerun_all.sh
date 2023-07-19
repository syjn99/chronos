#!/bin/sh
BASEDIR=/Users/jungwoo/Workspace/devnet/chronos/testnet/devnet
KAIROS_PATH=/Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/../../../kairos

# Kill chronos node if exists.

# change these to the unique parts of your command
chronos_part="p2p-tcp-port="

pids=$(ps aux | grep "${chronos_part}" | grep -v grep | awk '{print $2}')

if [ -z "$pids" ]
then
    echo "No Chronos processes found"
else
    echo "Killing Chronos processes with PIDs: $pids"
    for pid in $pids
    do
        kill -9 $pid
    done
fi


# Kill validator clients if exists.
validator_part="wallet-dir=$BASEDIR/validator"

pids=$(ps aux | grep "${validator_part}" | grep -v grep | awk '{print $2}')

if [ -z "$pids" ]
then
    echo "No validator client processes found"
else
    echo "Killing Validator processes with PIDs: $pids"
    for pid in $pids
    do
        kill -9 $pid
    done
fi

if [ "$1" = "stop" ]; then
    exit
fi

# Reset and run chronos nodes and validator clients.
rm -rf /Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/logs/chronos*
rm -rf /Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/logs/vali*

rm -rf /Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/node/beaconchaindata
rm -rf /Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/node/metaData

nohup /Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/node/run_node.sh > logs/chronos.out &
echo "Chronos node started"
nohup /Users/jungwoo/Workspace/devnet/chronos/testnet/devnet/validator/run_validator.sh > logs/vali.out &
    echo "Validator client started"
