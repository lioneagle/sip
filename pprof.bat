@echo off
setlocal
if exist benchmark.bat goto ok
echo benchmark.bat must be run from its folder
goto end
: ok
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0
set OLDGOBIN=%GOBIN%
set GOBIN=%~dp0bin

if not exist .\test_temp\%1.test.exe goto end
if not exist .\test_temp\%1_cpu.prof goto end

go tool pprof .\test_temp\%1.test.exe .\test_temp\%1_cpu.prof


:end
echo finished