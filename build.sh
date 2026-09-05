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

echo ">> Compiling for Windows 64-bit (dist/app.exe)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app.exe .

# 移除旧的冗余产物 (若存在)
rm -f ${OUTPUT_DIR}/app_tray.exe

echo ">> Compiling for current OS (dist/app)..."
go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app .

chmod +x ${OUTPUT_DIR}/app || true

echo "=========================================================="
echo " Build successful! Single binaries created in dist/:"
echo "   - dist/app.exe      (Windows 统一全能二进制，支持 CLI 命令与托盘自动无黑框)"
echo "   - dist/app          (macOS/Linux 统一全能二进制)"
echo "=========================================================="
