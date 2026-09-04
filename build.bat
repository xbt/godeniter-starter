@echo off
REM ==============================================================================
REM Godeniter Starter Windows 构建脚本
REM ==============================================================================

if not exist dist mkdir dist

if exist app.ico (
    echo ^>^> [ICON] 动态检测到应用图标，正在通过纯标准库生成 Windows 资源文件...
    go run github.com/xbt/godeniter/cmd/rsrc -auto
)

echo ^>^> Building single executable for Windows CLI...
go build -ldflags="-s -w" -o dist\app.exe .

echo ^>^> Building silent GUI tray client for Windows...
go build -ldflags="-s -w -H=windowsgui" -o dist\app_tray.exe .

if %ERRORLEVEL% equ 0 (
    echo Build successful! Outputs:
    echo   - dist\app.exe (Console mode)
    echo   - dist\app_tray.exe (Silent tray mode)
) else (
    echo Build failed with error %ERRORLEVEL%
)

pause
