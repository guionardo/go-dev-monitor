package config

import (
	"fmt"
	"os"
	"path"
	"sync"

	pathtools "github.com/guionardo/go-dev-monitor/internal/utils/path_tools"
)

var (
	getConfigDir = sync.OnceValue(func() string {
		cfgDir, err := tryGetConfigDir()
		if err != nil {
			fmt.Printf("CRITICAL: %s", err)
			os.Exit(2)
		}
		return cfgDir
	})

	getConfigDirFunc = os.UserConfigDir
)

const ConfigDirName = "go_dev_monitor"

func tryGetConfigDir() (string, error) {
	usrCfgDir, err := getConfigDirFunc()
	if err != nil {
		return "", fmt.Errorf("couldn't get config directory: %w", err)
	}
	cfgDir := path.Join(usrCfgDir, ConfigDirName)
	if err := pathtools.CreatePath(cfgDir); err != nil {
		return "", fmt.Errorf("couldn't create config dir: %s - %w", cfgDir, err)
	}
	if err := os.WriteFile(path.Join(cfgDir, "_write_test"), []byte{}, 0644); err != nil {
		return "", fmt.Errorf("couldn't write to config dir: %s - %w", cfgDir, err)
	}
	os.Remove(path.Join(cfgDir, "_write_test"))

	return cfgDir, nil
}

func GetDefaultConfigDir() (string, error) {
	usrCfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(usrCfgDir, ConfigDirName), nil
}

// SetConfigDir for tests
func SetConfigDir(configDir string) {
	getConfigDirFunc = func() (string, error) {
		if stat, err := os.Stat(configDir); err == nil && stat.IsDir() {
			return configDir, nil
		}
		return "", fmt.Errorf("directory not found: %s", configDir)
	}
}
