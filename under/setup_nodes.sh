BASEDIR=$(dirname "$0")
GETHPATH=$BASEDIR/kairos/build/bin/geth

# Clear former data
rm -rf $BASEDIR/nodes-*

# Create the shell scripts for each validator
for i in $(seq 0 1); do
    # Define the name of the new shell script
    script_name="$BASEDIR/node-$i/run_node.sh"

    # Calculate the port value based on the index
    httpport=$((22000 + i))
    wsport=$((32000 + i))
    tcpport=$((30303 + i))
    authport=$((8551 + i))

    # Consensus node ports
    rpcport=$((4000 + i))
    monitorport=$((8080 + i))
    udpV5port=$((12000 + i))
    tcpV5port=$((13000 + i))
    rpcgatewayport=$((3500 + i))

    # Copy the necessary files to the validator directories
    mkdir -p $BASEDIR/node-$i/keystore
    cp $BASEDIR/kairos/under/artifacts/keyfiles/keyfile$i.json $BASEDIR/nodes/node-$i/keystore/keyfile.json
    cp $BASEDIR/kairos/under/artifacts/nodekeys/nodekey$i $BASEDIR/nodes/node-$i/nodekey
    cp $BASEDIR/kairos/under/artifacts/nodekeys/nodekey$i.pub $BASEDIR/nodes/node-$i/nodekey.pub
    cp $BASEDIR/kairos/under/artifacts/accountPassword $BASEDIR/nodes/node-$i/accountPassword

    # Create the new shell script
    echo "#!/bin/sh" > "$script_name"

    # Add the provided code to the new shell script
    cat << EOF >> "$script_name"
BASEDIR=\$(dirname "\$0") 

bazel run //cmd/beacon-chain:beacon-chain -- \\
-datadir=\$BASEDIR/node-$i \\
    --min-sync-peers=0 \\
    --chain-config-file=../.config.yml \\
    --config-file=../../config.yml \\
    --chain-id=813 \\
    --execution-endpoint=http://localhost:${authport} \\
    --accept-terms-of-use \\
    --jwt-secret=geth/jwtsecret \\
    --contract-deployment-block=0 \\
    --p2p-udp-port"=${udpV5port}" \\
    --p2p-tcp-port"=${tcpV5port}" \\
    --rpc-port"=${rpcport}" \\
    --monitoring-port"=${monitorport}" \\
    --grpc-gateway-port"=${rpcgatewayport}" \\
    --p2p-local-ip 127.0.0.1 \\
    --bootstrap-node=enr:-MK4QFRaGA2AU27anJtkKjaMvmQuxCxy3z_uTFhXHNnMXBHrersXSPhu49tpOL-jXJcZBsXswvagy1ukVNHw2naDRBeGAYgzHZlih2F0dG5ldHOIAAAAAAAAAACEZXRoMpCMkQYoIAAAkgD2AgAAAAAAgmlkgnY0gmlwhMCoATWJc2VjcDI1NmsxoQOKuxPm3Kv03lZW9fWRHQ0XI7NMVUUvPWuiwHDAbpCaI4hzeW5jbmV0cwCDdGNwgjLIg3VkcIIu4
EOF

    # Make the new shell script executable
    chmod +x "$script_name"

    $GETHPATH --datadir $BASEDIR/nodes/node-$i init $BASEDIR/genesis.json

done