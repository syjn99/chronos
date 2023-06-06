@echo off
SET BASEDIR=%~dp0

@REM bazel run //cmd/validator:validator wallet create -- --wallet-dir=$BASEDIR/wallet --keymanager-kind=imported --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt
%BASEDIR%..\bin\validator.exe wallet create --wallet-dir=%BASEDIR%wallet2 --keymanager-kind=imported --wallet-password-file=%BASEDIR%artifacts/wallet/password.txt

@REM bazel run //cmd/validator:validator accounts import -- --wallet-dir=$BASEDIR/wallet --keys-dir=$BASEDIR/artifacts/keyfiles --wallet-password-file=$BASEDIR/artifacts/wallet/password.txt
%BASEDIR%..\bin\validator.exe accounts import --wallet-dir=%BASEDIR%%wallet --keys-dir=%BASEDIR%artifacts/keyfiles --wallet-password-file=%BASEDIR%artifacts/wallet/password.txt
