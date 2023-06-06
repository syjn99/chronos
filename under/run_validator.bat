@echo off
SET BASEDIR=%~dp0

@REM bazel run --config=minimal //cmd/validator:validator -- --wallet-dir=$BASEDIR/wallet --chain-config-file=$BASEDIR/config.yml --config-file=$BASEDIR/config.yml --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt --accept-terms-of-use --force-clear-db
%BASEDIR%..\validator.exe --wallet-dir=%BASEDIR%wallet --chain-config-file=%BASEDIR%config.yml --config-file=%BASEDIR%config.yml --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt --accept-terms-of-use --force-clear-db
