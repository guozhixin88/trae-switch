#!/bin/bash
# Trae Switch 一键停止脚本

sudo pkill -9 -f "trae-switch.app" 2>/dev/null

if ps aux | grep -v grep | grep -q "trae-switch.app/Contents/MacOS/trae-switch"; then
    echo "❌ 停止失败"
else
    echo "✅ Trae Switch 已停止"
fi
