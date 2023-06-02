BASEDIR=$(pwd)

bazel run --config=minimal //cmd/validator:validator -- --wallet-dir=$BASEDIR/wallet --chain-config-file=$BASEDIR/config.yml --config-file=$BASEDIR/config.yml --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt --accept-terms-of-use --force-clear-db
