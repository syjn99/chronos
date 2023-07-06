BASEDIR=$(dirname "$0")

$BIN_PATH/beacon-chain \ 
    --datadir $BASEDIR/chronos-beacon-chain \
    --genesis-state $BASEDIR/chronos-beacon-chain/genesis.ssz \
    --chain-config-file $BASEDIR/chronos-beacon-chain/config.yml \
    --config-file $BASEDIR/chronos-beacon-chain/config.yml \
    --chain-id 813 \ 
    --min-sync-peers 0 \
    --execution-endpoint http://localhost:8551 \
    --accept-terms-of-use \
    --jwt-secret $BASEDIR/kairos/geth/jwtsecret \
    --contract-deployment-block 0 \
    --p2p-udp-port 12000 \
    --p2p-tcp-port 13000 \
    --rpc-port 4000 \
    --monitoring-port 8080 \
    --grpc-gateway-port 3500 \
    --p2p-local-ip 127.0.0.1 \ 
    --bootstrap-node $BASEDIR/chronos-beacon-chain/bootnode.yaml \
    --verbosity debug
