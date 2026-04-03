package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"trae-switch/internal/cert"
	"trae-switch/internal/config"
	"trae-switch/internal/hosts"
	"trae-switch/internal/proxy"
	"trae-switch/internal/truststore"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const pidFile = "/tmp/trae-switch.pid"

type App struct {
	ctx          context.Context
	certManager  *cert.CertificateManager
	hostsManager *hosts.HostsManager
	trustManager *truststore.TrustStoreManager
	proxyServer  *proxy.ProxyServer
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir, err := os.UserConfigDir()
	if err != nil {
		dataDir = os.Getenv("APPDATA")
	}
	if dataDir != "" {
		dataDir = filepath.Join(dataDir, "trae-switch")
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Printf("Failed to create data directory: %v", err)
		}
	}

	a.certManager = cert.NewCertificateManager(dataDir)
	a.hostsManager = hosts.NewHostsManager()
	a.proxyServer = proxy.NewProxyServer("127.0.0.1", 443)

	if err := a.certManager.LoadOrGenerateCA(); err != nil {
		log.Printf("Failed to load/generate CA: %v", err)
	}

	a.trustManager = truststore.NewTrustStoreManager(a.certManager.GetCACertPath())

	if _, err := config.Load(); err != nil {
		log.Printf("Failed to load config: %v", err)
	}

	log.Println("Application started successfully")

	// 写入 PID 文件
	writePidFile()

	// 启动一个 goroutine 来检查并自动启动代理
	go a.autoStartProxy()
}

// writePidFile 写入当前进程 PID 到文件
func writePidFile() {
	pid := os.Getpid()
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		log.Printf("Failed to write pid file: %v", err)
	}
}

// removePidFile 删除 PID 文件
func removePidFile() {
	if err := os.Remove(pidFile); err != nil {
		log.Printf("Failed to remove pid file: %v", err)
	}
}

// autoStartProxy 检查配置并自动启动代理
func (a *App) autoStartProxy() {
	// 等待 3 秒让 UI 完全加载
	select {
	case <-a.ctx.Done():
		return
	case <-time.After(3 * time.Second):
	}

	// 检查是否已经启动
	if a.IsProxyRunning() {
		log.Println("Proxy already running, skipping auto-start")
		return
	}

	// 检查配置是否完成
	isSet, _ := a.hostsManager.IsSet()
	isInstalled, _ := a.trustManager.IsInstalled()
	providers := config.GetProviders()

	log.Printf("Auto-start check: hostsSet=%v, certInstalled=%v, providers=%d", isSet, isInstalled, len(providers))

	// 如果配置已完成且有服务商，自动启动
	if isSet && isInstalled && len(providers) > 0 {
		log.Println("Auto-starting proxy...")

		// 生成服务器证书
		if err := a.certManager.GenerateServerCert("api.openai.com"); err != nil {
			log.Printf("Failed to generate server cert: %v", err)
			runtime.EventsEmit(a.ctx, "auto-start-error", "生成证书失败："+err.Error())
			return
		}

		a.proxyServer.SetCertificate(
			a.certManager.GetServerCertPEM(),
			a.certManager.GetServerKeyPEM(),
		)

		if err := a.proxyServer.Start(a.ctx); err != nil {
			log.Printf("Failed to auto-start proxy: %v", err)
			runtime.EventsEmit(a.ctx, "auto-start-error", "启动代理失败："+err.Error())
			return
		}

		log.Println("Proxy auto-started successfully")
		runtime.EventsEmit(a.ctx, "auto-start-success", "代理已自动启动")
	} else {
		log.Println("Auto-start skipped: configuration incomplete")
	}
}

func (a *App) GetStatus() map[string]interface{} {
	portAvailable, portProcess := proxy.CheckPortStatus(443)

	result := map[string]interface{}{
		"runningAsAdmin":   truststore.IsRunningAsAdmin(),
		"proxyRunning":     false,
		"proxyPort":        443,
		"hostsSet":         false,
		"certInstalled":    false,
		"portAvailable":    portAvailable,
		"portProcess":      portProcess,
		"activeProvider":   nil,
		"activeTargetURL":  "",
	}

	if a.hostsManager != nil {
		isSet, _ := a.hostsManager.IsSet()
		result["hostsSet"] = isSet
	}

	if a.trustManager != nil {
		isInstalled, _ := a.trustManager.IsInstalled()
		result["certInstalled"] = isInstalled
	}

	if a.proxyServer != nil {
		status := a.proxyServer.GetStatus()
		result["proxyRunning"] = status.Running
		result["activeTargetURL"] = status.TargetURL
	}

	provider := config.GetActiveProvider()
	if provider != nil {
		result["activeProvider"] = map[string]interface{}{
			"name":   provider.Name,
			"models": provider.Models,
		}
	}

	return result
}

func (a *App) IsRunningAsAdmin() bool {
	return truststore.IsRunningAsAdmin()
}

func (a *App) SetHosts() error {
	if !truststore.IsRunningAsAdmin() {
		return fmt.Errorf("需要管理员权限")
	}
	return a.hostsManager.Set()
}

func (a *App) RestoreHosts() error {
	if !truststore.IsRunningAsAdmin() {
		return fmt.Errorf("需要管理员权限")
	}
	return a.hostsManager.Restore()
}

func (a *App) IsHostsSet() bool {
	isSet, _ := a.hostsManager.IsSet()
	return isSet
}

func (a *App) InstallCertificate() error {
	if !truststore.IsRunningAsAdmin() {
		return fmt.Errorf("需要管理员权限")
	}

	if err := a.certManager.LoadOrGenerateCA(); err != nil {
		return fmt.Errorf("生成证书失败：%w", err)
	}

	return a.trustManager.Install()
}

func (a *App) UninstallCertificate() error {
	if !truststore.IsRunningAsAdmin() {
		return fmt.Errorf("需要管理员权限")
	}
	return a.trustManager.Uninstall()
}

func (a *App) IsCertificateInstalled() bool {
	isInstalled, _ := a.trustManager.IsInstalled()
	return isInstalled
}

func (a *App) StartProxy() error {
	if !truststore.IsRunningAsAdmin() {
		return fmt.Errorf("需要管理员权限监听 443 端口")
	}

	if err := a.certManager.GenerateServerCert("api.openai.com"); err != nil {
		return fmt.Errorf("生成服务器证书失败：%w", err)
	}

	a.proxyServer.SetCertificate(
		a.certManager.GetServerCertPEM(),
		a.certManager.GetServerKeyPEM(),
	)

	return a.proxyServer.Start(a.ctx)
}

func (a *App) StopProxy() error {
	return a.proxyServer.Stop()
}

func (a *App) IsProxyRunning() bool {
	return a.proxyServer != nil && a.proxyServer.IsRunning()
}

func (a *App) QuickStart() error {
	if !truststore.IsRunningAsAdmin() {
		return fmt.Errorf("需要管理员权限")
	}

	if !a.IsHostsSet() {
		if err := a.SetHosts(); err != nil {
			return fmt.Errorf("设置 hosts 失败：%w", err)
		}
	}

	if !a.IsCertificateInstalled() {
		if err := a.InstallCertificate(); err != nil {
			return fmt.Errorf("安装证书失败：%w", err)
		}
	}

	if !a.IsProxyRunning() {
		if err := a.StartProxy(); err != nil {
			return fmt.Errorf("启动代理失败：%w", err)
		}
	}

	return nil
}

func (a *App) QuickStop() error {
	if a.IsProxyRunning() {
		if err := a.StopProxy(); err != nil {
			log.Printf("停止代理失败：%v", err)
		}
	}

	if a.IsHostsSet() {
		if err := a.RestoreHosts(); err != nil {
			log.Printf("恢复 hosts 失败：%v", err)
		}
	}

	return nil
}

func (a *App) GetProviders() []map[string]interface{} {
	providers := config.GetProviders()
	result := make([]map[string]interface{}, 0, len(providers))
	for i, p := range providers {
		result = append(result, map[string]interface{}{
			"index":       i,
			"name":        p.Name,
			"openai_base": p.OpenAIBase,
			"models":      p.Models,
		})
	}
	return result
}

func (a *App) GetActiveProviderIndex() int {
	return config.GetActiveProviderIndex()
}

func (a *App) SetActiveProvider(index int) error {
	return config.SetActiveProvider(index)
}

func (a *App) AddProvider(name, openaiBase string, models []string) error {
	provider := config.Provider{
		Name:       name,
		OpenAIBase: openaiBase,
		Models:     models,
	}
	return config.AddProvider(provider)
}

func (a *App) UpdateProvider(index int, name, openaiBase string, models []string) error {
	provider := config.Provider{
		Name:       name,
		OpenAIBase: openaiBase,
		Models:     models,
	}
	return config.UpdateProvider(index, provider)
}

func (a *App) DeleteProvider(index int) error {
	return config.DeleteProvider(index)
}

func (a *App) shutdown(ctx context.Context) {
	log.Println("Application shutting down...")

	// 删除 PID 文件
	removePidFile()

	if a.IsProxyRunning() {
		if err := a.StopProxy(); err != nil {
			log.Printf("Failed to stop proxy: %v", err)
		}
	}
}
