@echo off
REM ==============================================================================
REM Godeniter Starter Windows 构建脚本
REM ==============================================================================

if not exist dist mkdir dist

if exist app.ico (
    echo ^>^> [ICON] 动态检测到应用图标，正在通过纯标准库生成 Windows 资源文件...
    go run github.com/xbt/godeniter/cmd/rsrc -auto
)

echo ^>^> Building single all-in-one executable for Windows (dist\app.exe)...
go build -ldflags="-s -w" -o dist\app.exe .

if exist dist\app_tray.exe del dist\app_tray.exe

if %ERRORLEVEL% equ 0 (
    echo Build successful! Outputs:
    echo   - dist\app.exe (Unified binary: CLI commands + silent tray auto-hide)
) else (
    echo Build failed with error %ERRORLEVEL%
)

pause
