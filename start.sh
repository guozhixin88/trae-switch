#!/bin/bash
# Trae Switch 启动脚本
# 保存为 ~/trae-switch/start.sh 然后双击运行

osascript -e 'do shell script "~/trae-switch/build/bin/trae-switch.app/Contents/MacOS/trae-switch" with administrator privileges'
