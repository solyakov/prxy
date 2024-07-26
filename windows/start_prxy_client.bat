REM Example directory structure:
REM   prxy-client/
REM     client.crt
REM     client.key
REM     ca.crt
REM     prxy-client.exe
REM     start_prxy_client.bat

cd /d %~dp0

REM This is where the client will listen for incoming browser connections.
REM Specify this as the HTTPS proxy in the browser.
set PRXY_LISTEN=localhost:8080

REM This is where the client will tunnel the browser connections to.
set PRXY_SERVER=1.2.3.4:443

REM This CA file is used to verify the server certificate.
REM The client and server certificates must be signed by the same CA.
set PRXY_CA=ca.crt

REM This key pair is used to authenticate the client to the server.
set PRXY_CERTIFICATE=client.crt
set PRXY_KEY=client.key

REM Tunnel timeout for idle connections and buffer size for data transfer.
set PRXY_TIMEOUT=60s
set PRXY_BUFFER=32768

start "" "prxy-client.exe"
exit
