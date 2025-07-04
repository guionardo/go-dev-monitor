package config

import (
	"log/slog"
	"os"
	"path"

	"github.com/guionardo/go-dev-monitor/internal/debug"
	pathtools "github.com/guionardo/go-dev-monitor/internal/utils/path_tools"
	"gopkg.in/yaml.v3"
)

type Configer interface {
	FileName() string
	Reset()
	SetConfigDir(string)
	GetRoots() []string
}

func getConfigFile(filename string) (string, error) {
	usrCfgDir, err := getConfigDirFunc()
	if err != nil {
		return "", err
	}
	cfgDir := path.Join(usrCfgDir, "go_dev_monitor")
	if err := pathtools.CreatePath(cfgDir); err != nil {
		return "", err
	}
	return path.Join(cfgDir, filename), nil
}

func ReadConfigFile[T Configer](cfg T) error {
	fileName := cfg.FileName()
	fileName, err := getConfigFile(fileName)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(fileName)
	if err == nil {
		if err = yaml.Unmarshal(content, cfg); err == nil {
			cfg.Reset()
		}
	} else if os.IsNotExist(err) {
		cfg.Reset()
		err = nil
	}
	cfg.SetConfigDir(path.Dir(fileName))
	debug.Log().Debug("config", slog.String("filename", fileName), slog.Any("config", cfg))

	return err
}

func SaveConfigFile[T Configer](cfg T) error {
	fileName, err := getConfigFile((cfg.FileName()))
	if err != nil {
		return err
	}
	content, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, content, 0644)
}
