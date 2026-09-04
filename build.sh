#!/bin/bash
# ==============================================================================
# Godeniter Starter 跨平台打包脚本
# ==============================================================================

set -e

OUTPUT_DIR="./dist"
mkdir -p ${OUTPUT_DIR}

echo ">> Compiling for Windows 64-bit (dist/app.exe)..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app.exe main.go

echo ">> Compiling for current OS (dist/app)..."
go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app main.go

chmod +x ${OUTPUT_DIR}/app || true

echo "=========================================================="
echo " Build successful! Single binary created in dist/:"
echo "   - dist/app.exe (Windows 双击运行)"
echo "   - dist/app (macOS/Linux)"
echo "=========================================================="
