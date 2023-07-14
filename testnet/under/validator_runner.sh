BASEDIR=$(pwd)
KAIROS_PATH=$BASEDIR/../../../kairos

if [ "$1" = "clean" ]; then
    # Clear former data
    rm -rf $BASEDIR/validator-*
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
    for i in $(seq $start $end); do
        rm -rf $BASEDIR/validator-$i
        mnemonic=$(yq e ".[$i].mnemonic" "$BASEDIR/artifacts/mnemonics.yaml")
        count=$(yq e ".[$i].count" "$BASEDIR/artifacts/mnemonics.yaml")

        echo count is $count

        mkdir -p $BASEDIR/validator-$i

        printf "%s" "$mnemonic" >> $BASEDIR/validator-$i/mnemonic.txt

        # Recover wallet from mnemonic
        bazel run //cmd/validator:validator wallet recover -- --wallet-dir=$BASEDIR/validator-$i --mnemonic-file=$BASEDIR/validator-$i/mnemonic.txt --mnemonic-25th-word-file=$BASEDIR/artifacts/wallet/password.txt --num-accounts=$count --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt --accept-terms-of-use

        rpcport=$((7000 + i))
        beaconrpcport=$((4000 + i))
        rpcgatewayport=$((3500 + i))
        slasherrpc=$((4002 + i))

        # Define the name of the new shell script
        script_name="$BASEDIR/validator-$i/run_validator.sh"

        # Create the new shell script
        echo "#!/bin/sh" > "$script_name"

        # Add the provided code to the new shell script
        cat << EOF >> "$script_name"
BASEDIR=$(pwd)

bazel run //cmd/validator:validator -- \\
    --wallet-dir=$BASEDIR/validator-$i \\
    --proposer-settings-file=$BASEDIR/artifacts/recipients/recipient$i.yaml \\
    --chain-config-file=$BASEDIR/artifacts/config.yml \\
    --config-file=$BASEDIR/artifacts/config.yml \\
    --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt \\
    --beacon-rpc-provider=127.0.0.1:${beaconrpcport} \\
    --beacon-rpc-gateway-provider=127.0.0.1:${rpcgatewayport} \\
    --beacon-rest-api-provider=http://127.0.0.1:${rpcgatewayport} \\
    --rpc-port=${rpcport} \\
    --slasher-rpc-provider=127.0.0.1:${slasherrpc} \\
    --accept-terms-of-use \\
    --force-clear-db
EOF

        # Make the new shell script executable
        chmod +x "$script_name"
    done


elif [ "$1" = "stop" ]; then
    # change these to the unique parts of your command
    unique_part="wallet-dir=$BASEDIR/validator-"

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
    rm -rf $BASEDIR/logs/vali-*
    mkdir $BASEDIR/logs

    for num in $(seq $start $end)
    do
            nohup $BASEDIR/validator-$num/run_validator.sh > logs/vali-$num.out &
    done
else
    echo "Invalid argument. should be one of below
    clean - clear validator datas
    init n1 (n2) - Make initialized validator data from 0 to n1 (or n1 to n2). Max value 1 => 2 validators.
    stop - stop running all validators
    run n1 (n2) - run validators from 0 to n1 (or n1 to n2). Max value 1 => 2 validators."
fi