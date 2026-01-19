package gitops

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetGlobalconfig 从全局git配置中返回用户名称和电子邮件。
func GetGlobalConfig() (string, string, error) {
	nameCmd := exec.Command("git", "config", "--global", "--get", "user.name")
	prepareCmd(nameCmd)
	nameBytes, _ := nameCmd.Output()

	emailCmd := exec.Command("git", "config", "--global", "--get", "user.email")
	prepareCmd(emailCmd)
	emailBytes, _ := emailCmd.Output()

	name := strings.TrimSpace(string(nameBytes))
	email := strings.TrimSpace(string(emailBytes))

	return name, email, nil
}

// 设置全局配置，包括用户名称、邮箱和核心SSH命令
func SetGlobalConfig(name, email, sshKeyPath string) error {
	// 使用 --replace-all 强制覆盖可能存在的多个值
	if err := runGit("config", "--global", "--replace-all", "user.name", name); err != nil {
		return fmt.Errorf("设置用户名失败: %w", err)
	}
	if err := runGit("config", "--global", "--replace-all", "user.email", email); err != nil {
		return fmt.Errorf("设置邮箱失败: %w", err)
	}

	if sshKeyPath != "" {
		// Windows路径标准化用于SSH命令
		sshPath := strings.ReplaceAll(sshKeyPath, "\\", "/")
		sshCmd := fmt.Sprintf("ssh -i \"%s\" -o IdentitiesOnly=yes", sshPath)
		if err := runGit("config", "--global", "--replace-all", "core.sshCommand", sshCmd); err != nil {
			return fmt.Errorf("设置 SSH Key 失败: %w", err)
		}
	} else {
		// 如果为空则取消设置以使用默认值，使用 --unset-all 确保清除所有旧值
		_ = runGit("config", "--global", "--unset-all", "core.sshCommand")
	}

	return nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	prepareCmd(cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", err, string(output))
	}
	return nil
}
