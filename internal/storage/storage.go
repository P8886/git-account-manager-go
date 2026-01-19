package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Account struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	SSHKeyPath string `json:"ssh_key_path"`
}

type Config struct {
	Accounts []Account `json:"accounts"`
	ActiveID string    `json:"active_id"`
}

func GetConfigDir() (string, error) {
	home, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "git-account-manager-go")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}
	return dir, nil
}

func LoadConfig() (*Config, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "accounts.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Config{Accounts: []Account{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(config *Config) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, "accounts.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
