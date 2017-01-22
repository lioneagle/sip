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

if not exist test_temp mkdir test_temp
if exist .\test_temp\%1.test.exe  del .\test_temp\%1.test.exe

go test %1 -bench=. -cpuprofile=.\test_temp\%1_cpu.prof %2
if not exist ./test_temp/%1_cpu.prof goto end

go test %1 -bench=. -c

if not exist %1.test.exe goto end
move %1.test.exe .\test_temp\%1.test.exe

:end
echo finished