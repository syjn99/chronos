BASEDIR=$(pwd)
KAIROS_PATH=$BASEDIR/../../../kairos

# Get the OS name
os_name=$(uname)

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

    # Clear former data
    rm -rf $BASEDIR/node-*

    # Replace Genesis timestamp for new beacon chain
    current_date=$(date +%s)
    target_date=$((current_date + 60))

    if [ "$os_name" = "Linux" ]; then
        echo "This machine is Linux machine."
        echo "Target genesis time updated to : $(date -d @$target_date)"
        ip_address=$(hostname -I | awk '{print $1}')
        # ip_address=$(curl ifconfig.me)
    elif [ "$os_name" = "Darwin" ]; then
        echo "This machine is macOS machine."
        echo "Target genesis time updated to : $(date -r $target_date)"
        ip_address=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1')
    else
        echo "This machine is neither Linux nor macOS. So there can be some problems."
    fi
    echo "IP of this machine is $ip_address"

    bazel run --config=minimal //tools/change-genesis-timestamp -- \
        -genesis-state=$BASEDIR/artifacts/genesis.ssz \
        -timestamp=$target_date

    # Create the shell scripts for each validator
    for i in $(seq $start $end); do
        mkdir -p $BASEDIR/node-$i
        cp $BASEDIR/artifacts/network_keys/network-keys$i $BASEDIR/node-$i/network-keys

        # Define the name of the new shell script
        script_name="$BASEDIR/node-$i/run_node.sh"

        # Calculate the port value based on the index
        authport=$((8651 + i))
        rpcport=$((4100 + i))
        monitorport=$((9080 + i))
        udpport=$((14000 + i))
        tcpport=$((15000 + i))
        rpcgatewayport=$((4500 + i))

        # Copy the necessary files to the validator directories
        mkdir -p $BASEDIR/node-$i

        # Create the new shell script
        echo "#!/bin/sh" > "$script_name"

        # Add the provided code to the new shell script
        cat << EOF >> "$script_name"
KAIROS_PATH=$KAIROS_PATH/testnet/under/node-$i/geth

bazel run --config=minimal //cmd/beacon-chain:beacon-chain -- \\
    -datadir=$BASEDIR/node-$i \\
    -genesis-state=$BASEDIR/artifacts/genesis.ssz \\
    -chain-config-file=$BASEDIR/artifacts/config.yml \\
    -config-file=$BASEDIR/artifacts/config.yml \\
    -chain-id=822 \\
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
    -p2p-host-ip"=${ip_address}" \\
    -bootstrap-node=$BASEDIR/artifacts/bootnode.yaml \\
    -verbosity=debug
EOF

        # Make the new shell script executable
        chmod +x "$script_name"

    done


elif [ "$1" = "stop" ]; then
    # change these to the unique parts of your command
    unique_part="chain-id=822"

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
    init n1 (n2) - Make initialized node data from 0 to n1 (or n1 to n2). Max value 1 => 2 nodes.
    stop - stop running all nodes
    run n1 (n2) - run nodes from 0 to n1 (or n1 to n2). Max value 1 => 2 nodes"
fi