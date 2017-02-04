@echo off
setlocal
if exist coverage.bat goto ok
echo coverage.bat must be run from its folder
goto end
: ok

call env.bat

if not exist test_temp mkdir test_temp

if exist .\test_temp\%1_coverage.out  del .\test_temp\%1_coverage.out
if exist .\test_temp\%1_coverage.html del .\test_temp\%1_coverage.html

go test -coverprofile=./test_temp/%1_coverage.out %1 %2
if not exist ./test_temp/%1_coverage.out goto end

go tool cover -html=./test_temp/%1_coverage.out -o ./test_temp/%1_coverage.html
.\test_temp\%1_coverage.html

:end
echo finished