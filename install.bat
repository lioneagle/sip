@echo off
setlocal
if exist install.bat goto ok
echo install.bat must be run from its folder
goto end
: ok
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
set OLDGOBIN=%GOBIN%
set GOBIN=%~dp0bin
gofmt -w src
go install %1
:end
echo finished