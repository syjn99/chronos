@echo off
SETLOCAL ENABLEDELAYEDEXPANSION

set BASEDIR=%~dp0
echo %BASEDIR%

SET PRYSMCTLPATH=%BASEDIR%..\prysmctl.exe
echo %PRYSMCTLPATH%


rem Clear former data
for /d %%i in ("%BASEDIR%node-*") do (
    echo hi
    echo %%i | findstr /r /c:"[0-9]*$">nul && (
        echo Deleting: %%i
        rmdir /s /q "%%i"
    )
)

rem Run prysmctl to generate genesis
"%PRYSMCTLPATH%" ^
    testnet generate-genesis ^
    --output-ssz="%BASEDIR%genesis.ssz" ^
    --chain-config-file="%BASEDIR%config.yml" ^
    --geth-genesis-json-in="%BASEDIR%..\..\..\kairos_window\under\artifacts\genesis.json" ^
    --deposit-json-file="%BASEDIR%artifacts\deposits\deposit_data_under.json" ^
    --num-validators=0 ^
    --execution-endpoint="http://localhost:22000" ^
    --override-eth1data="true"

for /L %%i in (0,1,1) do (
    mkdir "%BASEDIR%node-%%i" >nul 2>&1
    copy "%BASEDIR%artifacts\network_keys\network-keys%%i" "%BASEDIR%node-%%i\network-keys"

    rem Define the name of the new batch script
    set "script_name=%BASEDIR%node-%%i\run_node.bat"
    set "node_dir=%BASEDIR%node-%%i\"
    rem Calculate the port value based on the index
    set /a "authport=8551 + %%i"
    set /a "rpcport=4000 + %%i"
    set /a "monitorport=8080 + %%i"
    set /a "udpport=12000 + %%i"
    set /a "tcpport=13000 + %%i"
    set /a "rpcgatewayport=3500 + %%i"
    set "kairos_jwt_path=%BASEDIR%..\..\..\kairos_window\under\node-%%i\geth\jwtsecret"


    (   
        echo @echo off 
        echo SET "BASEDIR=!node_dir!"
        echo SET "CHRONOS_PATH=%BASEDIR%..\..\beacon-chain.exe"
        echo SET "kairos_jwt_path=!kairos_jwt_path!"
        echo SET "authport=!authport!"
        echo SET "rpcport=!rpcport!"
        echo SET "monitorport=!monitorport!"
        echo SET "udpport=!udpport!"
        echo SET "tcpport=!tcpport!"
        echo SET "rpcgatewayport=!rpcgatewayport!"
    ) >> !script_name!
    type "artifacts\run_node.bat" >> !script_name!
)

ENDLOCAL