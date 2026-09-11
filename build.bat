@echo off
REM ==============================================================================
REM Godeniter Starter Windows 构建脚本
REM ==============================================================================

if not exist dist mkdir dist

if exist app.ico (
    echo ^>^> [ICON] 动态检测到应用图标，正在通过纯标准库生成 Windows 资源文件...
    go run github.com/xbt/godeniter/cmd/rsrc -auto
)

echo ^>^> 正在编译 Windows 统一全能程序 (dist\app.exe，桌面双击无黑框 + 终端支持完整CLI)...
go build -ldflags="-s -w -H=windowsgui" -o dist\app.exe .

if exist dist\app_tray.exe del dist\app_tray.exe

if %ERRORLEVEL% equ 0 (
    echo.
    echo ==========================================================
    echo 构建成功！生成单一全能可执行程序：
    echo   - dist\app.exe
    echo.
    echo 特性说明：
    echo   1. 桌面双击：100%% 彻底无黑框、无闪现，直接静默直达屏幕右下角托盘！
    echo   2. 命令行使用：终端中支持 run/console/start/stop/status 等全部指令输出！
    echo ==========================================================
) else (
    echo Build failed with error %ERRORLEVEL%
)

pause
