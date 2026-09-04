@echo off
REM ==============================================================================
REM Godeniter Starter Windows 构建脚本
REM ==============================================================================

if not exist dist mkdir dist

if exist app.ico (
    echo ^>^> [ICON] 动态检测到应用图标，正在通过纯标准库生成 Windows 资源文件...
    go run github.com/xbt/godeniter/cmd/rsrc -auto
)

echo ^>^> Building single executable for Windows...
go build -ldflags="-s -w" -o dist\app.exe .

if %ERRORLEVEL% equ 0 (
    echo Build successful! Output: dist\app.exe
) else (
    echo Build failed with error %ERRORLEVEL%
)

pause
