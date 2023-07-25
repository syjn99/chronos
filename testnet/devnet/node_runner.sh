BASEDIR=$(pwd)
KAIROS_PATH=$BASEDIR/../../../kairos

# Get the OS name
os_name=$(uname)

minimal=""

authport=8551
rpcport=4000
monitorport=8080
udpport=12000
tcpport=13000
rpcgatewayport=3500

if [ "$1" = "clean" ]; then
    # Clear former data
    cd -
    rm -rf $BASEDIR/node-*
    exit 0

elif [ "$1" = "init" ]; then
    # Update over-devnet repository to recent one
    cd $BASEDIR/../../../over-devnet
    git fetch origin
    git checkout master
    git pull

    # copy genesis & config for main devnet
    cp ./under_devnet/files/genesis.ssz ../chronos/testnet/devnet/genesis.ssz
    cp ./under_devnet/files/config.yml ../chronos/testnet/devnet/config.yml
    cp ./under_devnet/files/bootnode.yaml ../chronos/testnet/devnet/bootnode.yaml
    cd -
    
elif [ "$1" = "pver" ]; then
    echo "minimal build"
    minimal="--config=minimal "

    # Update over-devnet repository to recent one
    cd $BASEDIR/../../../over-devnet
    git fetch origin
    git checkout master
    git pull

    # copy genesis & config for minimal devnet (pver)
    cp ./pver_devnet/files/genesis_minimal.ssz ../chronos/testnet/devnet/genesis.ssz
    cp ./pver_devnet/files/config_minimal.yml ../chronos/testnet/devnet/config.yml
    cp ./pver_devnet/files/bootnode.yaml ../chronos/testnet/devnet/bootnode.yaml
    cd -

elif [ "$1" = "stop" ]; then

    # change these to the unique parts of your command
    unique_part="p2p-tcp-port=$tcpport"

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
    rm -rf $BASEDIR/logs/chronos*

    nohup $BASEDIR/node/run_node.sh > $BASEDIR/logs/chronos.out &
    exit 0
else
    echo "Invalid argument. should be one of below
    clean - clear node data
    init - Make initialized beacon node data for main devnet (many validator network)
    pver - Make initialized beacon node data for pver devnet (minimal validator network)
    stop - stop running all nodes
    run  - run node"
    exit 0
fi

# Set ip address for beacon node
if [ "$os_name" = "Linux" ]; then
    echo "This machine is Linux machine."
    ip_address=$(hostname -I | awk '{print $1}')
    # uncomment to use public IP for host-ip flag
    # ip_address=$(curl ifconfig.me)
elif [ "$os_name" = "Darwin" ]; then
    echo "This machine is macOS machine."
    ip_address=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1')
else
    echo "This machine is neither Linux nor macOS. So there can be some problems."
fi
echo "IP of this machine is $ip_address"

# Clear former data
rm -rf $BASEDIR/node

# Create the shell script for beacon node
mkdir -p $BASEDIR/node

# Define the name of the new shell script
script_name="$BASEDIR/node/run_node.sh"

# Create the new shell script
echo "#!/bin/sh" > "$script_name"

# Add the provided code to the new shell script
cat << EOF >> "$script_name"
KAIROS_PATH=$KAIROS_PATH/testnet/devnet/node-0/geth

bazel run $minimal//cmd/beacon-chain:beacon-chain -- \\
-datadir=$BASEDIR/node \\
-genesis-state=$BASEDIR/genesis.ssz \\
-chain-config-file=$BASEDIR/config.yml \\
-config-file=$BASEDIR/config.yml \\
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
-bootstrap-node=$BASEDIR/bootnode.yaml \\
-enable-upnp \\
-verbosity=debug
EOF

# Make the new shell script executable
chmod +x "$script_name"

## Make rerun script
# Clear former data
rm -rf $BASEDIR/rerun.sh
# Define the name of the new shell script
rerun_script="$BASEDIR/rerun_all.sh"
# Create the new shell script
    echo "#!/bin/sh" > "$rerun_script"

# Add the provided code to the new shell script
cat << EOF >> "$rerun_script"
BASEDIR=$(pwd)
KAIROS_PATH=$BASEDIR/../../../kairos

# Kill chronos node if exists.

# change these to the unique parts of your command
chronos_part="p2p-tcp-port=$tcpport_base"

pids=\$(ps aux | grep "\${chronos_part}" | grep -v grep | awk '{print \$2}')

if [ -z "\$pids" ]
then
    echo "No Chronos processes found"
else
    echo "Killing Chronos processes with PIDs: \$pids"
    for pid in \$pids
    do
        kill -9 \$pid
    done
fi


# Kill validator clients if exists.
validator_part="wallet-dir=\$BASEDIR/validator"

pids=\$(ps aux | grep "\${validator_part}" | grep -v grep | awk '{print \$2}')

if [ -z "\$pids" ]
then
    echo "No validator client processes found"
else
    echo "Killing Validator processes with PIDs: \$pids"
    for pid in \$pids
    do
        kill -9 \$pid
    done
fi

if [ "\$1" = "stop" ]; then
    exit
fi

# Reset and run chronos nodes and validator clients.
rm -rf $BASEDIR/logs/chronos*
rm -rf $BASEDIR/logs/vali*

rm -rf $BASEDIR/node/beaconchaindata
rm -rf $BASEDIR/node/metaData

nohup $BASEDIR/node/run_node.sh > logs/chronos.out &
echo "Chronos node started"
nohup $BASEDIR/validator/run_validator.sh > logs/vali.out &
    echo "Validator client started"
EOF

# Make the new shell script executable
chmod +x "$rerun_script"