BASEDIR=$(pwd)
echo $BASEDIR

# Clear former data
rm -rf $BASEDIR/node-*
rm -rf $BASEDIR/bootnode.yaml

# Run the command and capture the output
# bazel run //tools/enr-calculator:enr-calculator -- --private 534a9f6de7c84cea0ef5d04e86c3ff7616843cb5f2a820a29ef175dada89f2c6 --ipAddress 127.0.0.1 --udp-port 12000 --tcp-port 13000 --out $BASEDIR/bootnode.yaml
enr-calculator.exe -- --private 534a9f6de7c84cea0ef5d04e86c3ff7616843cb5f2a820a29ef175dada89f2c6 --ipAddress 127.0.0.1 --udp-port 12000 --tcp-port 13000 --out $BASEDIR/bootnode.yaml
# enr-calculator.exe -- --private 534a9f6de7c84cea0ef5d04e86c3ff7616843cb5f2a820a29ef175dada89f2c6 --ipAddress 127.0.0.1 --udp-port 12000 --tcp-port 13000 --out %BASEDIR%\bootnode.yaml


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
KAIROS_PATH=$BASEDIR/../../kairos/under/node-$i/geth
echo \$KAIROS_PATH

beacon-chain.exe -- \\
    -datadir=$BASEDIR/node-$i \\
    -min-sync-peers=0 \\
    -chain-config-file=$BASEDIR/config.yml \\
    -config-file=$BASEDIR/config.yml \\
    -chain-id=813 \\
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