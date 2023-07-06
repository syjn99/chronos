BASEDIR=$(dirname "$0")

$BIN_PATH/validator \ 
    --wallet-dir $BASEDIR/chronos-validator \
    --proposer-setting-file $BASEDIR/chronos-validator/recipient.yaml \
    --chain-config=file $BASEDIR/chronos-validator/config.yml \
    --wallet-password-file $BASEDIR/chronos-validator/password.txt \
    --beacon-rpc-provider http://localhost:4000 \
    --beacon-rpc-gateway-provider http://localhost:3500 \
    --beacon-rest-api-provider http://localhost:8080 \
    --rpc-port 4001 \
    --slasher-rpc-provider http://localhost:4000 \
    --accept-terms-of-use \
    --force-clear-db
