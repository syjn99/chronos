@echo off
SET BASEDIR=%~dp0

@REM bazel run --config=minimal //cmd/validator:validator -- --wallet-dir=$BASEDIR/wallet --chain-config-file=$BASEDIR/config.yml --config-file=$BASEDIR/config.yml --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt --accept-terms-of-use --force-clear-db
@REM "%BASEDIR%..\validator.exe" --wallet-dir="%BASEDIR%wallet" --chain-config-file="%BASEDIR%config.yml" --config-file="%BASEDIR%config.yml" --wallet-password-file="%BASEDIR%artifacts/wallet/password.txt" --accept-terms-of-use --force-clear-db


"%VALIDATOR_CLIENT_PATH%" ^
    --wallet-dir="%BASEDIR%" ^
    --proposer-settings-file="%BASEDIR%recipient.yaml" ^
    --chain-config-file="%BASEDIR%config.yml" ^
    --wallet-password-file="%BASEDIR%password.txt" ^
    --beacon-rpc-provider="127.0.0.1:%beaconrpcport% ^
    --beacon-rpc-gateway-provider=127.0.0.1:%rpcgatewayport% ^
    --beacon-rest-api-provider=http://127.0.0.1:%rpcgatewayport% ^
    --rpc-port=%rpcport% ^
    --slasher-rpc-provider=127.0.0.1:%slasherrpc% ^
    --accept-terms-of-use ^
    --force-clear-db 
