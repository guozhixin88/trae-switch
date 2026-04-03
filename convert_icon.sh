#!/bin/bash
# 图标转换脚本

# 检查是否安装了图标转换工具
if ! which iconutil > /dev/null 2>&1; then
    echo "iconutil not found"
    exit 1
fi

# 创建 iconset 目录
mkdir -p ~/trae-switch/icon.iconset

# 复制原始图标（假设是 PNG 格式）
cp ~/trae-switch/build/appicon.png ~/trae-switch/icon.iconset/icon_512x512@2x.png

# 生成其他尺寸的图标
sips -z 16 16 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_16x16.png
sips -z 32 32 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_16x16@2x.png
sips -z 32 32 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_32x32.png
sips -z 64 64 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_32x32@2x.png
sips -z 128 128 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_128x128.png
sips -z 256 256 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_128x128@2x.png
sips -z 256 256 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_256x256.png
sips -z 512 512 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_256x256@2x.png
sips -z 512 512 ~/trae-switch/icon.iconset/icon_512x512@2x.png --out ~/trae-switch/icon.iconset/icon_512x512.png

# 转换为 icns 格式
iconutil -c icns ~/trae-switch/icon.iconset -o ~/trae-switch/TraeSwitchLauncher.app/Contents/Resources/icon.icns

# 清理
rm -rf ~/trae-switch/icon.iconset

echo "图标转换完成！"
