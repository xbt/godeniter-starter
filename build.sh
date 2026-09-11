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
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o ${OUTPUT_DIR}/app.exe .

rm -f ${OUTPUT_DIR}/app_tray.exe

echo ">> Compiling for current OS (dist/app)..."
go build -ldflags="-s -w" -o ${OUTPUT_DIR}/app .

chmod +x ${OUTPUT_DIR}/app || true

# 如果在 macOS 系统上，自动组装无需终端弹窗的原生桌面应用包 (dist/Godeniter.app)
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo ">> Packaging native macOS Status Bar App (dist/Godeniter.app)..."
    APP_BUNDLE="${OUTPUT_DIR}/Godeniter.app"
    mkdir -p "${APP_BUNDLE}/Contents/MacOS"
    mkdir -p "${APP_BUNDLE}/Contents/Resources"
    cp "${OUTPUT_DIR}/app" "${APP_BUNDLE}/Contents/MacOS/app"
    chmod +x "${APP_BUNDLE}/Contents/MacOS/app"

    # 若有应用图标，自动生成 macOS 原生完整标准 10 阶 .icns 高清图标资源
    if [ -f "app.ico" ]; then
        TEMP_ICONSET="/tmp/AppIcon.iconset"
        rm -rf "$TEMP_ICONSET" /tmp/base_icon.png
        mkdir -p "$TEMP_ICONSET"
        sips -s format png app.ico --out /tmp/base_icon.png >/dev/null 2>&1 || true

        sips -z 16 16     /tmp/base_icon.png --out "$TEMP_ICONSET/icon_16x16.png" >/dev/null 2>&1 || true
        sips -z 32 32     /tmp/base_icon.png --out "$TEMP_ICONSET/icon_16x16@2x.png" >/dev/null 2>&1 || true
        sips -z 32 32     /tmp/base_icon.png --out "$TEMP_ICONSET/icon_32x32.png" >/dev/null 2>&1 || true
        sips -z 64 64     /tmp/base_icon.png --out "$TEMP_ICONSET/icon_32x32@2x.png" >/dev/null 2>&1 || true
        sips -z 128 128   /tmp/base_icon.png --out "$TEMP_ICONSET/icon_128x128.png" >/dev/null 2>&1 || true
        sips -z 256 256   /tmp/base_icon.png --out "$TEMP_ICONSET/icon_128x128@2x.png" >/dev/null 2>&1 || true
        sips -z 256 256   /tmp/base_icon.png --out "$TEMP_ICONSET/icon_256x256.png" >/dev/null 2>&1 || true
        sips -z 512 512   /tmp/base_icon.png --out "$TEMP_ICONSET/icon_256x256@2x.png" >/dev/null 2>&1 || true
        sips -z 512 512   /tmp/base_icon.png --out "$TEMP_ICONSET/icon_512x512.png" >/dev/null 2>&1 || true
        sips -z 1024 1024 /tmp/base_icon.png --out "$TEMP_ICONSET/icon_512x512@2x.png" >/dev/null 2>&1 || true

        iconutil -c icns "$TEMP_ICONSET" -o "${APP_BUNDLE}/Contents/Resources/AppIcon.icns" >/dev/null 2>&1 || true
        rm -rf "$TEMP_ICONSET" /tmp/base_icon.png
    fi

    cat << 'EOF' > "${APP_BUNDLE}/Contents/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>app</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.godeniter.starter</string>
    <key>CFBundleName</key>
    <string>Godeniter</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.5</string>
    <key>LSUIElement</key>
    <false/>
</dict>
</plist>
EOF

    # 强制刷新 macOS LaunchServices 与 Finder 对该 Bundle 的图标缓存
    touch "${APP_BUNDLE}"
    /System/Library/Frameworks/CoreServices.framework/Versions/A/Frameworks/LaunchServices.framework/Versions/A/Support/lsregister -f "${APP_BUNDLE}" >/dev/null 2>&1 || true
    killall Finder >/dev/null 2>&1 || true
fi

echo "=========================================================="
echo " Build successful! Single binaries created in dist/:"
echo "   - dist/app.exe        (Windows 统一全能二进制，桌面双击无黑框直达托盘，终端支持完整CLI)"
echo "   - dist/app            (macOS/Linux 统一全能二进制，终端执行直接在顶部状态栏常驻)"
if [[ "$OSTYPE" == "darwin"* ]]; then
echo "   - dist/Godeniter.app  (macOS 原生状态栏应用，访达双击无黑框直达顶部菜单栏)"
fi
echo "=========================================================="
