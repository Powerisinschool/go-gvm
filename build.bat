@ECHO OFF
REM "Builds an installable executable for GVM"
md bin
SET GOBIN=%CD%\bin

echo ----------------------------
echo Building gvm.exe
echo ----------------------------
cd .\src
go build gvm.go

cd ..\
move .\src\gvm.exe %GOBIN%

echo Consider adding the following path to your environmental variables:
echo %GOBIN%