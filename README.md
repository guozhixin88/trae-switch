# Trae Switch

**Trae Switch** 是一个专为 Trae IDE 设计的工具，通过 DNS 劫持 + 本地反向代理，让 Trae IDE 支持第三方大模型服务商 API（如阿里百炼 Coding Plan、kimi coding plan 等）。详细使用教程：https://mp.weixin.qq.com/s/W_Z_nbrO7ioU8upcq4KkOw

## 🚀 功能特点

- **多服务商支持**：可添加、编辑、删除多个服务商配置
- **本地模型管理**：`/v1/models` 请求返回本地配置的模型（因为第三方通常不支持此接口）
- **自动 Hosts 配置**：将 `api.openai.com` 重定向到 `127.0.0.1`
- **CA 证书管理**：生成并安装本地 CA 证书，用于 HTTPS 拦截
- **实时状态监控**：显示代理运行状态和当前激活的服务商
- **不需要输入 key**：通过在 trae 中配置 key，在本工具不需要输入任何 apikey

## 📋 支持的服务商

- ✅ 阿里百炼、Kimi 等 Coding Plan
- ✅ 其他支持 OpenAI 协议的第三方 api 服务商

## 🔧 技术架构

### 技术栈
- **后端**：Go (Wails 框架)
- **前端**：Svelte + Tailwind CSS
- **网络**：HTTPS 代理服务器
- **系统**：跨平台支持（Windows/macOS）

### 核心模块
1. **代理服务器**：监听 443 端口，处理 OpenAI API 请求
2. **配置管理**：读写 `config.json` 配置文件
3. **Hosts 管理**：自动设置和恢复 Hosts 配置
4. **证书管理**：生成和安装自签名 CA 证书
5. **前端界面**：现代化的用户交互界面

## 📦 安装方法

### macOS 用户

#### 前置要求
- Go 1.22+ 
- Node.js 16+
- Xcode Command Line Tools

#### 方法一：从源码构建（推荐）
1. **安装 Go**（如果未安装）：
   ```bash
   brew install go
   ```

2. **克隆仓库**：
   ```bash
   git clone https://github.com/guozhixin88/trae-switch.git
   cd trae-switch
   ```

3. **安装 Wails CLI**：
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **安装前端依赖**：
   ```bash
   cd frontend
   npm install
   cd ..
   ```

5. **构建应用**：
   ```bash
   wails build
   ```
   构建完成后，应用位于 `build/bin/trae-switch.app`

6. **运行应用**（需要管理员权限）：
   ```bash
   sudo ./build/bin/trae-switch.app/Contents/MacOS/trae-switch
   ```
   或
   ```bash
   sudo wails dev
   ```

### Windows 用户

#### 方法一：直接下载可执行文件
1. 从 [Releases](https://github.com/yourusername/trae-switch/releases) 页面下载最新版本的 `trae-switch.exe`
2. 以管理员身份运行程序

#### 方法二：从源码构建
1. 确保已安装 Go 1.22+ 和 Node.js 16+
2. 克隆仓库：
   ```bash
   git clone https://github.com/yourusername/trae-switch.git
   cd trae-switch
   ```
3. 安装 Wails CLI：`go install github.com/wailsapp/wails/v2/cmd/wails@latest`
4. 安装前端依赖：`cd frontend && npm install`
5. 构建应用：`wails build`
6. 运行：`./build/bin/trae-switch.exe`

## 🛠️ 使用方法

### 1. 添加服务商配置
1. 点击「+ 添加」按钮
2. 填写服务商信息：
   - **名称**：服务商显示名称（如 "阿里百炼"）
   - **API 地址**：OpenAI 协议的 API 地址（如 `https://coding.dashscope.aliyuncs.com/v1`）
   - **模型列表**：添加可用的模型 ID（如 `qwen3.5-plus`、`kimi-k2.5` 等）
3. 点击「保存」

### 2. 启动代理
1. 确保系统配置中的「Hosts 配置」和「CA 证书」都已启用
2. 点击右上角的「启动」按钮
3. 代理启动成功后，状态栏会显示「运行中」

### 3. 在 Trae IDE 中使用
1. 打开 Trae IDE
2. 进入模型设置
3. 选择 OpenAI 服务商
4. 输入对应第三方服务商的真实 API Key（如 `sk-xxx`）
5. 手动输入你想要使用的模型名称
6. 关闭 auto mode 并选择刚添加的模型
7. 开始使用！

---

## 📖 快速开始指南（以阿里云百炼为例）

### 第一步：获取 API Key
1. 访问 [阿里云百炼控制台](https://bailian.console.aliyun.com/)
2. 登录阿里云账号
3. 进入 **API-KEY 管理** 页面
4. 创建新的 API Key（格式：`sk-sp-xxxxx`）
5. ⚠️ **复制并妥善保存**（只显示一次）

### 第二步：配置 Trae Switch
1. 启动 Trae Switch 应用
2. 点击「+ 添加」添加服务商
3. 填写以下信息：
   - **名称**：`阿里云百炼`
   - **API 地址**：`https://coding.dashscope.aliyuncs.com/v1`
   - **模型列表**：
     - `qwen3.5-plus`
     - `qwen3-coder-plus`
     - `qwen3-coder-next`
     - `kimi-k2.5`
     - `glm-5`
     - （或其他你需要的模型）
4. 点击「保存」
5. 确保「Hosts 配置」和「CA 证书」显示已启用
6. 点击「启动」按钮

### 第三步：在 Trae IDE 中配置
1. 打开 Trae IDE
2. 进入 **Settings** → **AI Model**（或类似设置）
3. 添加新模型：
   - **Provider（服务商）**：选择 **OpenAI**
   - **Base URL**：`https://api.openai.com/v1`（保持默认，不要修改！）
   - **API Key**：输入第一步获取的阿里云 API Key
   - **Model Name**：输入 `qwen3.5-plus`
4. 保存配置

### 第四步：开始使用
1. 在 Trae IDE 聊天窗口中
2. 关闭 **Auto Mode**
3. 选择你刚添加的 `qwen3.5-plus` 模型
4. 发送消息测试

---

## ⚠️ 重要注意事项

### macOS 用户必读

1. **443 端口权限**
   - 应用需要监听 443 端口，必须使用 `sudo` 运行
   - 如果提示端口被占用，请关闭占用 443 端口的程序

2. **VPN 冲突问题**
   - 某些 VPN 可能会劫持 DNS 或修改路由表
   - 如遇连接问题，请尝试暂时关闭 VPN
   - 或者在 VPN 设置中将 `api.openai.com` 添加到白名单/不走 VPN

3. **系统缓存**
   - 如遇 DNS 解析问题，可执行以下命令清除缓存：
     ```bash
     sudo dscacheutil -flushcache
     sudo killall -HUP mDNSResponder
     ```

4. **证书信任**
   - 首次运行时，系统可能会弹出证书警告
   - 这是正常的，因为应用需要安装自签名证书
   - 在应用内点击「安装证书」即可

## ⚙️ 配置文件

配置文件位于程序同目录下的 `config.json`，格式如下：

```json
{
  "providers": [
    {
      "name": "服务商名称",
      "openai_base": "https://api.example.com/v1",
      "models": ["model-1", "model-2"]
    }
  ],
  "active_provider": 0
}
```

- `name`：服务商显示名称
- `openai_base`：OpenAI 协议的 API 地址
- `models`：模型 ID 列表
- `active_provider`：当前激活的服务商索引

## 📝 使用说明

1. 添加服务商配置（API 地址和模型列表）并点击「启动」
2. 在 Trae IDE 添加自定义模型，服务商选择 OpenAI 服务商
3. 模型手动输入你想要使用的模型并且输入对应 API Key（如 sk-xxx）
4. 关闭 auto mode 并且选择刚添加的模型

## 🔍 常见问题

### Q: 启动失败怎么办？
**A:** 请检查：
- 是否以管理员身份运行（macOS 需要使用 sudo）
- 443 端口是否被占用
- Hosts 配置是否成功
- CA 证书是否安装

### Q: 如何检查 443 端口是否被占用？
**A:** 执行以下命令：
```bash
sudo lsof -i :443 -sTCP:LISTEN
```
如果有其他程序占用，请关闭该程序后再试。

### Q: macOS 上如何卸载证书？
**A:** 运行以下命令：
```bash
sudo security delete-certificate -c "Trae Switch Root CA" /Library/Keychains/System.keychain
```

### Q: macOS 上如何恢复 hosts？
**A:** 程序退出时会自动恢复，也可以手动编辑 `/etc/hosts` 删除 Trae Switch 相关条目：
```bash
sudo nano /etc/hosts
# 删除以下行：
# === Trae Switch Start ===
127.0.0.1 api.openai.com
# === Trae Switch End ===
```

### Q: 模型不显示怎么办？
**A:** 请确保：
- 已在服务商配置中添加了模型
- 已选择了正确的服务商
- 代理已成功启动
- **关闭 VPN 后重试**（VPN 可能会劫持 DNS）

### Q: 提示 "Incorrect model name" 错误
**A:** 请检查：
1. Trae Switch 是否正常运行
2. hosts 是否正确：`grep api.openai.com /etc/hosts` 应返回 `127.0.0.1 api.openai.com`
3. 在 Trae IDE 中，Base URL 必须是 `https://api.openai.com/v1`（不是阿里云地址！）
4. 尝试清除 DNS 缓存：
   ```bash
   sudo dscacheutil -flushcache
   sudo killall -HUP mDNSResponder
   ```

### Q: 如何验证代理是否正常工作？
**A:** 执行以下命令测试：
```bash
curl -k https://127.0.0.1/v1/models
```
如果返回包含模型列表的 JSON，说明代理正常。

### Q: API Key 如何获取？
**A:** API Key 需要从对应服务商的官方网站获取：
- 阿里云百炼：https://bailian.console.aliyun.com/
- Kimi：https://platform.moonshot.cn/

### Q: 支持哪些模型？
**A:** 支持所有支持 OpenAI 接口协议服务商提供的模型，只要在配置中添加对应的模型 ID 即可。

### Q: 如何配置多个服务商？
**A:** 可以添加多个服务商配置，通过 `active_provider` 字段切换当前使用的服务商。在应用界面中点击对应的服务商即可激活。

### Q: VPN 会影响使用吗？
**A:** 会！某些 VPN 会：
- 修改系统 DNS 设置
- 劫持所有 HTTPS 流量
- 导致 `api.openai.com` 的请求不经过本地代理

**解决方案**：
1. 暂时关闭 VPN
2. 或在 VPN 设置中将 `api.openai.com` 添加到排除列表

## 🛡️ 安全性

- **本地运行**：所有数据处理都在本地进行，不会上传任何数据
- **自签名证书**：仅用于本地 HTTPS 拦截，不会影响其他应用
- **Hosts 修改**：仅修改 `api.openai.com` 的解析，不影响其他域名
- **不存储 key**：通过在 trae 中配置 key，在本工具不需要输入任何 key


## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

---
## Star History

<a href="https://www.star-history.com/?repos=mtfly%2Ftrae-switch&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/image?repos=mtfly/trae-switch&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/image?repos=mtfly/trae-switch&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/image?repos=mtfly/trae-switch&type=date&legend=top-left" />
 </picture>
</a>
