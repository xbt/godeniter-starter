@echo off
REM ==============================================================================
REM Godeniter Starter Windows 构建脚本
REM ==============================================================================

if not exist dist mkdir dist

if exist app.ico (
    echo ^>^> [ICON] 动态检测到应用图标，正在通过纯标准库生成 Windows 资源文件...
    go run github.com/xbt/godeniter/cmd/rsrc -auto
)

echo ^>^> 正在编译 Windows 统一全能程序 (dist\app.exe)...
go build -ldflags="-s -w" -o dist\app.exe .

echo ^>^> 正在编译 Windows 纯静默托盘客户端 (dist\app_tray.exe，彻底无黑框)...
go build -ldflags="-s -w -H=windowsgui" -o dist\app_tray.exe .

if %ERRORLEVEL% equ 0 (
    echo.
    echo ==========================================================
    echo 构建成功！生成产物说明：
    echo   - dist\app_tray.exe : 纯桌面托盘客户端（双击直接进入右下角托盘，彻底无黑框）
    echo   - dist\app.exe      : 统一全能二进制（双击进入托盘自动隐藏黑框；终端支持 run/start/stop 命令）
    echo ==========================================================
) else (
    echo Build failed with error %ERRORLEVEL%
)

pause
