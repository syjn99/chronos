BASEDIR=$(pwd)
KAIROS_PATH=$BASEDIR/../../../kairos

# Get the OS name
os_name=$(uname)
minimal=""

rpcport=8000
beaconrpcport=4000
rpcgatewayport=3500
slasherrpc=4002

if [ "$1" = "clean" ]; then
    # Clear former data
    rm -rf $BASEDIR/validator
    exit 0

elif [ "$1" = "init" ]; then
    echo "normal build"

elif [ "$1" = "pver" ]; then
    echo "minimal build"
    minimal="--config=minimal "

elif [ "$1" = "stop" ]; then
    # change these to the unique parts of your command
    unique_part="wallet-dir=$BASEDIR/validator"

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
    exit 0

elif [ "$1" = "run" ]; then
    rm -rf $BASEDIR/logs/vali*
    mkdir $BASEDIR/logs
    nohup $BASEDIR/validator/run_validator.sh > logs/validator.out &
    exit 0

else
    echo "Invalid argument. should be one of below
    clean - clear validator datas
    init - Make initialized validator data for main devnet (many validator network)
    pver - Make initialized validator data for pver devnet (minimal validator network)
    stop - stop running validator
    run - run validator"
    exit 0

fi

rm -rf $BASEDIR/validator

mnemonic=$(yq e ".[1].mnemonic" "$BASEDIR/artifacts/mnemonics.yaml")
count=$(yq e ".[1].count" "$BASEDIR/artifacts/mnemonics.yaml")
echo count is $count

mkdir -p $BASEDIR/validator
printf "%s" "$mnemonic" >> $BASEDIR/validator/mnemonic.txt

# Recover wallet from mnemonic
bazel run $minimal//cmd/validator:validator wallet recover -- --wallet-dir=$BASEDIR/validator --mnemonic-file=$BASEDIR/validator/mnemonic.txt --mnemonic-25th-word-file=$BASEDIR/artifacts/wallet/password.txt --num-accounts=$count --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt --accept-terms-of-use

# Define the name of the new shell script
script_name="$BASEDIR/validator/run_validator.sh"

# Create the new shell script
echo "#!/bin/sh" > "$script_name"

# Add the provided code to the new shell script
cat << EOF >> "$script_name"
BASEDIR=$(pwd)

bazel run $minimal//cmd/validator:validator -- \\
--wallet-dir=$BASEDIR/validator \\
--proposer-settings-file=$BASEDIR/artifacts/recipient.yaml \\
--chain-config-file=$BASEDIR/config.yml \\
--config-file=$BASEDIR/config.yml \\
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