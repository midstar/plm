@echo off
REM Builds plm using the correct version, build time and git hash
if [%1] EQU [] (
  echo Usage:
  echo     install_plm.bat 'version'
  exit /b 1
)
echo version: %1
for /f %%i in ('git rev-parse HEAD') do set GITHASH=%%i
echo git hash: %GITHASH%
echo building / installing
set INSTALLCMD=go install -ldflags="-X main.applicationBuildTime=%DATE%_%TIME% -X main.applicationVersion=%1 -X main.applicationGitHash=%GITHASH%" github.com/midstar/plm
echo %INSTALLCMD%
%INSTALLCMD%
