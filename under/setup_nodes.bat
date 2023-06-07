@echo off
SETLOCAL ENABLEDELAYEDEXPANSION

set BASEDIR=%~dp0
echo %BASEDIR%

SET ENRPATH=%BASEDIR%..\enr-calculator.exe
echo %ENRPATH%
rem Clear former data
for /d %%i in ("%BASEDIR%node-*") do (
    echo hi
    echo %%i | findstr /r /c:"[0-9]*$">nul && (
        echo Deleting: %%i
        rmdir /s /q "%%i"
    )
)

del /S /Q "%BASEDIR%bootnode.yaml" >nul 2>&1

rem Run the command and capture the output
"%ENRPATH%" ^
    --private 534a9f6de7c84cea0ef5d04e86c3ff7616843cb5f2a820a29ef175dada89f2c6 ^
    --ipAddress 127.0.0.1 ^
    --udp-port 12000 ^
    --tcp-port 13000 ^
    --out "%BASEDIR%bootnode.yaml"


for /L %%i in (0,1,1) do (
    mkdir "%BASEDIR%\node-%%i" >nul 2>&1
    copy "%BASEDIR%\artifacts\network_keys\network-keys%%i" "%BASEDIR%\node-%%i\network-keys"

    rem Define the name of the new batch script
    set "script_name=%BASEDIR%\node-%%i\run_node.bat"
    set "node_dir=%BASEDIR%\node-%%i\"
    rem Calculate the port value based on the index
    set /a "authport=8551 + %%i"
    set /a "rpcport=4000 + %%i"
    set /a "monitorport=8080 + %%i"
    set /a "udpport=12000 + %%i"
    set /a "tcpport=13000 + %%i"
    set /a "rpcgatewayport=3500 + %%i"
    set "kairos_jwt_path=%BASEDIR%..\..\kairos_window\under\node-%%i\geth\jwtsecret"

    rem Copy the necessary files to the validator directories
    mkdir "%BASEDIR%\node-%%i" >nul 2>&1

    (   
        echo @echo off 
        @REM echo SET "BASEDIR2=%%~dp0"
        echo SET "BASEDIR=!node_dir!"
        echo SET "CHRONOS_PATH=%BASEDIR%..\beacon-chain.exe"
        echo SET "kairos_jwt_path=!kairos_jwt_path!"
        echo SET "authport=!authport%!"
        echo SET "rpcport=!rpcport%!"
        echo SET "monitorport=!monitorport!"
        echo SET "udpport=!udpport!"
        echo SET "tcpport=!tcpport!"
        echo SET "rpcgatewayport=!rpcgatewayport!"
    ) >> !script_name!
    type "artifacts\run_node.bat" >> !script_name!
)

ENDLOCAL