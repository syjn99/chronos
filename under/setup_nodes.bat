@echo off
set BASEDIR=%cd%
echo %BASEDIR%

rem Clear former data
del /S /Q %BASEDIR%\node-* >nul 2>&1
del /S /Q %BASEDIR%\bootnode.yaml >nul 2>&1

rem Run the command and capture the output
enr-calculator.exe -- --private 534a9f6de7c84cea0ef5d04e86c3ff7616843cb5f2a820a29ef175dada89f2c6 --ipAddress 127.0.0.1 --udp-port 12000 --tcp-port 13000 --out %BASEDIR%\bootnode.yaml

rem Create the batch scripts for each validator
for /L %%i in (0,1,1) do (
    mkdir %BASEDIR%\node-%%i >nul 2>&1
    copy %BASEDIR%\artifacts\network_keys\network-keys%%i %BASEDIR%\node-%%i\network-keys

    rem Define the name of the new batch script
    set "script_name=%BASEDIR%\node-%%i\run_node.cmd"

    rem Calculate the port value based on the index
    set /a "authport=8551 + %%i"
    set /a "rpcport=4000 + %%i"
    set /a "monitorport=8080 + %%i"
    set /a "udpport=12000 + %%i"
    set /a "tcpport=13000 + %%i"
    set /a "rpcgatewayport=3500 + %%i"

    rem Copy the necessary files to the validator directories
    mkdir %BASEDIR%\node-%%i >nul 2>&1

    rem Create the new batch script
    (
        echo @echo off
        echo set "KAIROS_PATH=%BASEDIR%\..\..\kairos\under\node-%%i\geth"
        echo echo ^%KAIROS_PATH^%
        echo beacon-chain.exe -- ^
        echo    -datadir=%BASEDIR%\node-%%i ^
        echo    -min-sync-peers=0 ^
        echo    -chain-config-file=%BASEDIR%\config.yml ^
        echo    -config-file=%BASEDIR%\config.yml ^
        echo    -chain-id=813 ^
        echo    -execution-endpoint=http://localhost:^%authport^% ^
        echo    -accept-terms-of-use ^
        echo    -jwt-secret=^%KAIROS_PATH^%\jwtsecret ^
        echo    -contract-deployment-block=0 ^
        echo    -p2p-udp-port="^%udpport^%" ^
        echo    -p2p-tcp-port="^%tcpport^%" ^
        echo    -rpc-port="^%rpcport^%" ^
        echo    -monitoring-port="^%monitorport^%" ^
        echo    -grpc-gateway-port="^%rpcgatewayport^%" ^
        echo    -p2p-local-ip 127.0.0.1 ^
        echo    -bootstrap-node=%BASEDIR%\bootnode.yaml ^
        echo    -verbosity=debug
    ) > "%script_name%"
)
