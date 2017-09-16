@echo off
setlocal
if exist benchmark.bat goto ok
echo benchmark.bat must be run from its folder
goto end
: ok
call env.bat

if not exist .\test_temp\%1.test.exe goto end
if not exist .\test_temp\%1_cpu.prof goto end
gofmt -w src

if "%2" == "" (
	go tool pprof .\test_temp\%1.test.exe .\test_temp\%1_cpu.prof
) else (
	echo go tool pprof %2 -nodecount=%3 .\test_temp\%1.test.exe .\test_temp\%1_cpu.prof
	go tool pprof %2 -nodecount=%3 .\test_temp\%1.test.exe .\test_temp\%1_cpu.prof
)


:end
echo finished