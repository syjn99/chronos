BASEDIR=$(pwd)

for i in $(seq 0 1); do
    rm -rf $BASEDIR/validator-$i
    bazel run --config=minimal //cmd/validator:validator wallet create -- --wallet-dir=$BASEDIR/validator-$i --keymanager-kind=imported --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt
    bazel run --config=minimal //cmd/validator:validator accounts import -- --wallet-dir=$BASEDIR/validator-$i --keys-dir=$BASEDIR/artifacts/keyfiles/validator$i --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt

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

bazel run --config=minimal //cmd/validator:validator -- \\
    --wallet-dir=$BASEDIR/validator-$i \\
    --proposer-settings-file=$BASEDIR/artifacts/recipients/recipient$i.yaml \\
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
done