package config

import (
	"log/slog"
	"os"
	"path"

	"github.com/guionardo/go-dev-monitor/internal/debug"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent      *Agent  `yaml:"agent,omitempty"`
	Server     *Server `yaml:"server,omitempty"`
	configFile string
}

func NewConfig() *Config {
	configFile := path.Join(getConfigDir(), "config.yaml")
	logger := debug.Log().With(slog.String("filename", configFile))

	content, err := os.ReadFile(configFile)
	cfg := &Config{configFile: configFile}
	if err == nil {
		err = yaml.Unmarshal(content, cfg)
	}
	if err == nil {
		logger.Debug("config file")
	} else {
		logger.Warn("config file", slog.Any("error", err))
	}
	return cfg
}

func (cfg *Config) Save() error {
	logger := debug.Log().With(slog.String("filename", cfg.configFile))
	content, err := yaml.Marshal(cfg)
	if err == nil {
		err = os.WriteFile(cfg.configFile, content, 0644)
	}
	if err == nil {
		logger.Debug("saved config")
	} else {
		logger.Error("saving config", slog.Any("error", err))
	}
	return err
}

func (cfg *Config) GetConfigDir() string {
	return path.Dir(cfg.configFile)
}
