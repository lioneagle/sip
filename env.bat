set OLDGOPATH=%GOPATH%
set GOPATH=%OLDGOPATH%;%~dp0
set OLDGOBIN=%GOBIN%
set GOBIN=%~dp0bin
rem set GODEBUG=gctrace=1 2 > log_file 