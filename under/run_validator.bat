set BASEDIR=%cd%

validator.exe -- --wallet-dir=%BASEDIR%\wallet --chain-config-file=%BASEDIR%\config.yml --config-file=%BASEDIR%\config.yml --wallet-password-file=%BASEDIR%\artifacts\wallet\password.txt
