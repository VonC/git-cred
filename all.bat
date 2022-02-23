@echo off
setlocal enabledelayedexpansion

for %%i in ("%~dp0.") do SET "script_dir=%%~fi"
cd "%script_dir%"
for %%i in ("%~dp0.") do SET "dirname=%%~ni"

if "%1" == "amd" ( 
    set "barg=amd"
    shift
)
call build.bat %barg%
if errorlevel 1 (
    echo ERROR BUILD 1>&2
    exit /b 1
)
call run.bat %*