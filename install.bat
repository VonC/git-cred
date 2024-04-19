@echo off
setlocal

for %%i in ("%~dp0.") do SET "script_dir=%%~fi"
cd "%script_dir%" || echo "unable to cd to '%script_dir%'"&& exit /b 1
setlocal enabledelayedexpansion
for %%i in ("%~dp0.") do SET "dirname=%%~ni"

if not exist batcolors\echos.bat (
    echo Missing submodules
    echo Executing 'git submodule update --init'
    git submodule update --init
    if errorlevel 1 (
        echo "Submodules not properly initialized"
        exit /b 1
    )
)
call  "%script_dir%\batcolors\echos_macros.bat"

if exist "%script_dir%\senv.bat" (
    call "%script_dir%\senv.bat"
)

%_task% "Must copy '%dirname%.exe' to '%GOPATH%\bin\git-grep.exe'"
cp -f "%dirname%.exe" "%GOPATH%\bin\git-grep.exe"
if errorlevel 1 (
    %_fatal% "Unable to copy '%dirname%.exe' to '%GOPATH%\bin\git-grep.exe'" 2
)
%_ok% "Successfully copied '%dirname%.exe' to '%GOPATH%\bin\git-grep.exe'"
