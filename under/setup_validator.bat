set BASEDIR=%cd%

validator.exe wallet create -- --wallet-dir=%BASEDIR%\wallet --keymanager-kind=imported --wallet-password-file=%BASEDIR%\artifacts\wallet\password.txt
validator.exe accounts import -- --wallet-dir=%BASEDIR%\wallet --keys-dir=%BASEDIR%\artifacts\keyfiles --wallet-password-file=%BASEDIR%\artifacts\wallet\password.txt
