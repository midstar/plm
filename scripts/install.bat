@echo off
rem run this script as admin

if not exist %GOPATH%\bin\plm.exe (
    echo Install plm before installing by running "go install github.com/midstar/plm"
    exit /b 1
)

sc create plm binpath= "%GOPATH%\bin\plm.exe %GOPATH%\src\github.com\midstar\plm" start= auto DisplayName= "Process Load Monitor"
sc description plm "Process Load Monitor Service"
sc start plm
sc query plm

echo Check plm