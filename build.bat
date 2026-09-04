@echo off
REM ==============================================================================
REM Godeniter Starter Windows 构建脚本
REM ==============================================================================

if not exist dist mkdir dist

echo ^>^> Building single executable for Windows...
go build -ldflags="-s -w" -o dist\app.exe main.go

if %ERRORLEVEL% equ 0 (
    echo Build successful! Output: dist\app.exe
) else (
    echo Build failed with error %ERRORLEVEL%
)

pause
