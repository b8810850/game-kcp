@echo off
set GOPATH=%GOPATH%;%~dp0
set GOROOT=C:\Go

set GO=C:\Go\bin\go.exe

set GOOS=windows
set GOARCH=amd64

echo start build PingKcpServer ...
%GO% build PingKcpServer.go

echo ok
pause