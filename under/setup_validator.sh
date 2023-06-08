BASEDIR=$(pwd)

bazel run //cmd/validator:validator wallet create -- --wallet-dir=$BASEDIR/wallet --keymanager-kind=imported --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt
bazel run //cmd/validator:validator accounts import -- --wallet-dir=$BASEDIR/wallet --keys-dir=$BASEDIR/artifacts/keyfiles --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt