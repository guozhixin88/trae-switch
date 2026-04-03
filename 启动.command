#!/bin/bash
# Trae Switch 一键启动脚本（无任何窗口）
# 使用方法：双击运行或放在桌面

# 检查是否已经在运行
if ps aux | grep -v grep | grep -q "trae-switch.app/Contents/MacOS/trae-switch"; then
    echo "Trae Switch 已经在运行中"
    exit 0
fi

# 停止可能存在的旧进程
sudo pkill -9 -f "trae-switch.app" 2>/dev/null

# 等待一秒
sleep 1

# 后台启动应用（无窗口）
sudo nohup ~/trae-switch/build/bin/trae-switch.app/Contents/MacOS/trae-switch > /tmp/trae-switch.log 2>&1 &

# 等待启动
sleep 2

# 验证是否启动成功
if ps aux | grep -v grep | grep -q "trae-switch.app/Contents/MacOS/trae-switch"; then
    echo "✅ Trae Switch 启动成功"
else
    echo "❌ 启动失败，请查看 /tmp/trae-switch.log"
fi
