@echo off
setlocal
if exist benchmark.bat goto ok
echo benchmark.bat must be run from its folder
goto end
: ok
call env.bat
gofmt -w src
if not exist test_temp mkdir test_temp
if exist .\test_temp\%1.test.exe  del .\test_temp\%1.test.exe

go test %1 -bench=%2 -cpuprofile=.\test_temp\%1_cpu.prof %3
if not exist ./test_temp/%1_cpu.prof goto end

go test %1 -bench=. -c

if not exist %1.test.exe goto end
move %1.test.exe .\test_temp\%1.test.exe

:end
echo finished