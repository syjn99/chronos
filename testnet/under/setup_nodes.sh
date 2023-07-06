BASEDIR=$(pwd)
echo $BASEDIR

# Clear former data
rm -rf $BASEDIR/node-*
# rm -rf $BASEDIR/bootnode.yaml

# Run the command and save bootnode.yaml
bazel run --config=minimal //cmd/prysmctl:prysmctl testnet generate-genesis -- \
    --output-ssz=$BASEDIR/genesis.ssz \
    --chain-config-file=$BASEDIR/config.yml \
    --geth-genesis-json-in=$BASEDIR/../../../kairos/testnet/under/artifacts/genesis.json \
    --deposit-json-file=$BASEDIR/artifacts/deposits/deposit_data_under.json \
    --num-validators=0 \
    --execution-endpoint=http://localhost:22000 \
    --override-eth1data=true

# Create the shell scripts for each validator
for i in $(seq 0 1); do
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
KAIROS_PATH=$BASEDIR/../../../kairos/under/node-$i/geth
echo \$KAIROS_PATH

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