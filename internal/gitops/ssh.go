package gitops

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ListSSHKeys 扫描默认 .ssh 目录下的私钥文件
func ListSSHKeys() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	sshDir := filepath.Join(home, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var keys []string
	excludedExts := map[string]bool{
		".pub": true,
		".ppk": true,
	}
	excludedNames := map[string]bool{
		"known_hosts":     true,
		"known_hosts.old": true,
		"config":          true,
		"authorized_keys": true,
		"environment":     true,
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if excludedNames[name] {
			continue
		}

		ext := strings.ToLower(filepath.Ext(name))
		if excludedExts[ext] {
			continue
		}

		fullPath := filepath.Join(sshDir, name)
		if isPrivateKey(fullPath) {
			keys = append(keys, fullPath)
		}
	}

	return keys, nil
}

// isPrivateKey 简单验证文件是否为私钥 (检查头部)
func isPrivateKey(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		return strings.Contains(line, "PRIVATE KEY") || strings.HasPrefix(line, "-----BEGIN")
	}
	return false
}

// GenerateSSHKey 生成新的 SSH 密钥
func GenerateSSHKey(name, email, keyType string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	sshDir := filepath.Join(home, ".ssh")
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		_ = os.Mkdir(sshDir, 0700)
	}

	keyPath := filepath.Join(sshDir, name)
	if _, err := os.Stat(keyPath); err == nil {
		return "", fmt.Errorf("密钥文件 '%s' 已存在", name)
	}

	// ssh-keygen -t <type> -C <email> -f <path> -N ""
	cmd := exec.Command("ssh-keygen", "-t", keyType, "-C", email, "-f", keyPath, "-N", "")
	prepareCmd(cmd) // 使用 cmd_windows.go/cmd_unix.go 中的 helper

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("生成密钥失败: %s\n%s", err, string(output))
	}

	return keyPath, nil
}
