//go:build darwin
// +build darwin

package truststore

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	ErrNotAdmin = errors.New("administrator privileges required")
)

type TrustStoreManager struct {
	certPath string
}

func NewTrustStoreManager(certPath string) *TrustStoreManager {
	return &TrustStoreManager{
		certPath: certPath,
	}
}

func (tm *TrustStoreManager) IsInstalled() (bool, error) {
	// Check if certificate exists in the System keychain
	cmd := exec.Command("security", "find-certificate", "-c", "Trae Switch Root CA", "/Library/Keychains/System.keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, nil // Certificate not found, not an error
	}
	return len(output) > 0 || err == nil, nil
}

func (tm *TrustStoreManager) Install() error {
	installed, err := tm.IsInstalled()
	if err != nil {
		return err
	}
	if installed {
		return nil
	}

	// Add certificate to System keychain and trust it
	// Using security add-trusted-cert with admin privileges
	cmd := exec.Command("security", "add-trusted-cert", "-d", "-k", "/Library/Keychains/System.keychain", tm.certPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "authorization failed") || strings.Contains(string(output), "Operation not permitted") {
			return fmt.Errorf("%w: %s", ErrNotAdmin, string(output))
		}
		return fmt.Errorf("failed to install certificate: %w, output: %s", err, string(output))
	}

	return nil
}

func (tm *TrustStoreManager) Uninstall() error {
	installed, err := tm.IsInstalled()
	if err != nil {
		return err
	}
	if !installed {
		return nil
	}

	// Delete the certificate from System keychain
	cmd := exec.Command("security", "delete-certificate", "-c", "Trae Switch Root CA", "/Library/Keychains/System.keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to uninstall certificate: %w, output: %s", err, string(output))
	}

	return nil
}

func IsRunningAsAdmin() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	// Root UID is 0
	return strings.TrimSpace(string(output)) == "0"
}
