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

for /f "delims=" %%i in ('type "%script_dir%\go.mod"') do (
    if "!module_name!" == "" (
        set "module_name=%%i"
        goto:fldone
    )
)
:fldone
set "module_name=%module_name:module =%"
echo module_name='%module_name%'
if "%module_name%" == "" (
        %_fatal% "go.mod in '%script_dir%' does not include a module name" 66
)

rem https://medium.com/@joshroppo/setting-go-1-5-variables-at-compile-time-for-versioning-5b30a965d33e
for /f %%i in ('git describe --long --tags --dirty --always') do set gitver=%%i
for /f %%i in ('git describe --tags 2^>NUL') do set VERSION=%%i
@echo off
if "%VERSION%"=="" ( set "VERSION=0.1.0-tag" )
if not "%VERSION:-=%" == "%VERSION%" (
    set "todelete=-%VERSION:*-=%"
    call set "VERSION=%%VERSION:!todelete!=%%"
    echo snap !VERSION! with todelete='!todelete!'
    set "patch=!VERSION:*.=!"
    set "patch=!patch:*.=!"
) else (
    echo release tag '%VERSION%'
    set "patch="
)
if not "%patch%" == "" (
    call set "preversion=%%VERSION:.!patch!=%%"
    rem no need to increment a tag patch version if there was no vX.Y.Z tag in the first place
    if not "%gitver:v=%" == "%gitver%" ( set /A patch=patch+1 )
)
if not "%patch%" == "" (
    set "VERSION=%preversion%.%patch%_snapshot"
)
@echo off
echo VERSION='%VERSION%', patch='%patch%', preversion='%preversion%'

rem https://superuser.com/questions/1287756/how-can-i-get-the-date-in-a-locale-independent-format-in-a-batch-file
rem https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/get-date?view=powershell-7.1
rem C:\Windows\System32\WindowsPowershell\v1.0\powershell -Command "Get-Date -format 'yyyy-MM-dd_HH-mm-ss K'"
%+@% for /f %%a in ('C:\Windows\System32\WindowsPowershell\v1.0\powershell -Command "Get-Date -format yyyy-MM-dd_HH-mm-ss"') do set dtStamp=%%a
rem SET dtStamp
echo "dtStamp='%dtStamp%'"

set outputname=%dirname%.exe

if "%1" == "amd" (
    set GOARCH=amd64
    set GOOS=linux
    set "outputname=%dirname%_%VERSION%"
    %_info% "AMD build requested for %module_name%"
    set "fflag=-gcflags="all=-N -l" "
    rem dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./exename
)

%_info% "Start Building"
go build %fflag%-ldflags "-X %module_name%/version.GitTag=%gitver% -X %module_name%/version.BuildUser=%USERNAME% -X %module_name%/version.Version=%VERSION% -X %module_name%/version.BuildDate=%dtStamp%" -o %outputname%

if errorlevel 1 (
    %_fatal% "ERROR BUILD %module_name%" 3
)
set "filenamee=test"
if "%1" == "amd" (
    if "%LINUX_ACCOUNT%" == "" (
        %_fatal% "LINUX_ACCOUNT environment variable must be set" 5
    )
    rem echo.filename2 before='%filenamee%'  ------------
    call:setfilename
    rem echo.filename AFTER='!filename!'
    %_info% "Start scp to %LINUX_ACCOUNT%:/home/%LINUX_ACCOUNT%/bin/!filename!"
    scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/nul -q !filename! %LINUX_ACCOUNT%:/home/%LINUX_ACCOUNT%/bin/!filename!
    if errorlevel 1 (
        %_fatal% "scp to %LINUX_ACCOUNT%:/home/%LINUX_ACCOUNT%/bin/!filename! failed" 6
    )
    ssh %LINUX_ACCOUNT% "chmod 755 /home/%LINUX_ACCOUNT%/bin/!filename!"
    rem ssh %LINUX_ACCOUNT% "ln -fs /home/%LINUX_ACCOUNT%/bin/!filename! /project/%LINUX_ACCOUNT%/bin/%module_name%""
)
goto:eof
rem if "%1" neq "" ( %dirname% %* )

:setfilename
:: Use WMIC to retrieve date and time
FOR /F "skip=1 tokens=1-6" %%G IN ('WMIC Path Win32_LocalTime Get Day^,Hour^,Minute^,Month^,Second^,Year /Format:table') DO (
   IF "%%~L"=="" goto s_done
      Set _yyyy=%%L
      Set _mm=00%%J
      Set _dd=00%%G
      Set _hour=00%%H
      SET _minute=00%%I
      SET _second=00%%K
)
:s_done

:: Pad digits with leading zeros
      Set _mm=%_mm:~-2%
      Set _dd=%_dd:~-2%
      Set _hour=%_hour:~-2%
      Set _minute=%_minute:~-2%
      Set _second=%_second:~-2%

rem set "filename=%outputname%_%_yyyy%%_mm%%_dd%-%_hour%%_minute%%_second%"
set "filename=%outputname%"
goto:eof