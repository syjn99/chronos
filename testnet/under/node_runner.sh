BASEDIR=$(pwd)
KAIROS_PATH=$BASEDIR/../../../kairos

if [ "$1" = "clean" ]; then
    # Clear former data
    rm -rf $BASEDIR/node-*
elif [ "$1" = "init" ]; then
    if [ $# -eq 2 ]; then
        if ! [[ $2 =~ [0-9]+$ ]]; then
            echo "Invalid argument. second argument for init should be number"
            exit 1
        fi
        start=0
        end=$2
    elif [ $# -eq 3 ]; then
        if ! [[ $2 =~ [0-9]+$ ]]; then
            echo "Invalid argument. second argument for init should be number"
            exit 1
        elif ! [[ $3 =~ [0-9]+$ ]]; then
            echo "Invalid argument. third argument for init should be number"
            exit 1
        fi
        start=$2
        end=$3
    fi

    # Run the command and save bootnode.yaml
    bazel run --config=minimal //cmd/prysmctl:prysmctl testnet generate-genesis -- \
        --output-ssz=$BASEDIR/genesis.ssz \
        --chain-config-file=$BASEDIR/config.yml \
        --geth-genesis-json-in=$KAIROS_PATH/testnet/under/artifacts/genesis.json \
        --deposit-json-file=$BASEDIR/artifacts/deposits/deposit_data_under.json \
        --num-validators=0 \
        --execution-endpoint=http://localhost:22000 \
        --override-eth1data=true

    # Create the shell scripts for each validator
    for i in $(seq $start $end); do
        mkdir -p $BASEDIR/node-$i
        cp $BASEDIR/artifacts/network_keys/network-keys$i $BASEDIR/node-$i/network-keys

        # Define the name of the new shell script
        script_name="$BASEDIR/node-$i/run_node.sh"

        # Calculate the port value based on the index
        authport=$((8551 + i))
        rpcport=$((4000 + i))
        monitorport=$((8080 + i))
        udpport=$((12000 + i))
        tcpport=$((13000 + i))
        rpcgatewayport=$((3500 + i))

        # Copy the necessary files to the validator directories
        mkdir -p $BASEDIR/node-$i

        # Create the new shell script
        echo "#!/bin/sh" > "$script_name"

        # Add the provided code to the new shell script
        cat << EOF >> "$script_name"
KAIROS_PATH=$KAIROS_PATH/testnet/under/node-$i/geth

bazel run --config=minimal //cmd/beacon-chain:beacon-chain -- \\
    -datadir=$BASEDIR/node-$i \\
    -genesis-state=$BASEDIR/genesis.ssz \\
    -chain-config-file=$BASEDIR/config.yml \\
    -config-file=$BASEDIR/config.yml \\
    -chain-id=813 \\
    -min-sync-peers=0 \\
    -execution-endpoint=http://localhost:${authport} \\
    -accept-terms-of-use \\
    -jwt-secret=\$KAIROS_PATH/jwtsecret \\
    -contract-deployment-block=0 \\
    -p2p-udp-port"=${udpport}" \\
    -p2p-tcp-port"=${tcpport}" \\
    -rpc-port"=${rpcport}" \\
    -monitoring-port"=${monitorport}" \\
    -grpc-gateway-port"=${rpcgatewayport}" \\
    -p2p-local-ip 127.0.0.1 \\
    -bootstrap-node=$BASEDIR/bootnode.yaml \\
    -verbosity=debug
EOF

        # Make the new shell script executable
        chmod +x "$script_name"

    done


elif [ "$1" = "stop" ]; then
    # change these to the unique parts of your command
    unique_part="chain-id=813"

    pids=$(ps aux | grep "${unique_part}" | grep -v grep | awk '{print $2}')

    if [ -z "$pids" ]
    then
        echo "No processes found with command parts $unique_part"
    else
        echo "Killing Chronos processes with PIDs: $pids"
        for pid in $pids
        do
            kill -9 $pid
        done
    fi
elif [ "$1" = "run" ]; then
    if [ $# -eq 2 ]; then
        if ! [[ $2 =~ [0-9]+$ ]]; then
            echo "Invalid argument. second argument for init should be number"
            exit 1
        fi
        start=0
        end=$2
    elif [ $# -eq 3 ]; then
        if ! [[ $2 =~ [0-9]+$ ]]; then
            echo "Invalid argument. second argument for init should be number"
            exit 1
        elif ! [[ $3 =~ [0-9]+$ ]]; then
            echo "Invalid argument. third argument for init should be number"
            exit 1
        fi
        start=$2
        end=$3
    fi
    rm -rf $BASEDIR/logs/chronos-*
    mkdir $BASEDIR/logs

    for num in $(seq $start $end)
    do
            nohup $BASEDIR/node-$num/run_node.sh > logs/chronos-$num.out &
    done
else
    echo "Invalid argument. should be one of below
    clean - clear node data
    init n1 (n2) - Make initialized node data from 0 to n1 (or n1 to n2)
    stop - stop running nodes
    run n1 (n2) - run nodes from 0 to n1 (or n1 to n2)"
fi