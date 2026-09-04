#!/bin/bash
# ==============================================================================
# Godeniter Starter 跨平台打包脚本
# ==============================================================================

set -e

OUTPUT_DIR="./dist"
mkdir -p ${OUTPUT_DIR}

# 动态检测是否存在 app.ico 图标，若存在则使用内置纯 Go 标准库工具自动生成 Windows 资源段
if [ -f "app.ico" ] || [ -f "favicon.ico" ]; then
    echo ">> [ICON] 动态检测到应用图标，正在通过纯标准库生成 Windows 资源文件 (resource_windows_amd64.syso)..."
    go run github.com/xbt/godeniter/cmd/rsrc -auto || true
fi

echo ">> Compiling for Windows 64-bit CLI (dist/app.exe)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app.exe .

echo ">> Compiling for Windows 64-bit Silent GUI Tray (dist/app_tray.exe)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o ${OUTPUT_DIR}/app_tray.exe .

echo ">> Compiling for current OS (dist/app)..."
go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app .

chmod +x ${OUTPUT_DIR}/app || true

echo "=========================================================="
echo " Build successful! Single binaries created in dist/:"
echo "   - dist/app.exe      (Windows 控制台守护进程/双击运行)"
echo "   - dist/app_tray.exe (Windows 纯静默托盘客户端，无黑框)"
echo "   - dist/app          (macOS/Linux 统一二进制)"
echo "=========================================================="
