@echo off
setlocal
if exist coverage.bat goto ok
echo coverage.bat must be run from its folder
goto end
: ok

call env.bat

if not exist test_temp mkdir test_temp

if exist .\test_temp\coverage.out  del .\test_temp\coverage.out
if exist .\test_temp\coverage.html del .\test_temp\coverage.html

go test -coverprofile=./test_temp/coverage.out %1 %2
if not exist ./test_temp/coverage.out goto end

go tool cover -html=./test_temp/coverage.out -o ./test_temp/coverage.html
.\test_temp\coverage.html

:end
echo finished