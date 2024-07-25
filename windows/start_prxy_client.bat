@echo off
set PRXY_SERVER=localhost:8081
set PRXY_CERTIFICATE=C:\path\to\client.crt
set PRXY_KEY=C:\path\to\client.key
set PRXY_CA=C:\path\to\ca.crt
set PRXY_LISTEN=localhost:8080
set PRXY_TIMEOUT=60s
set PRXY_BUFFER=32768
start /min "" "C:\path\to\prxy-client.exe"
exit
