@echo off
setlocal enabledelayedexpansion

for %%i in ("%~dp0.") do SET "script_dir=%%~fi"
cd "%script_dir%"
for %%i in ("%~dp0.") do SET "dirname=%%~ni"

if exist senv.bat (
    call senv.bat
)
rem @echo on
if not "%PAGER_LOG%" == "" (
    del "%PAGER_LOG%" 2>NUL
)
call "%dirname%.exe" %*
if errorlevel 1 (
    if not "%PAGER_LOG%" == "" (
        type "%PAGER_LOG%" 2>NUL
    )
)