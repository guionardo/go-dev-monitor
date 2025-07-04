package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	pathtools "github.com/guionardo/go-dev-monitor/internal/utils/path_tools"
	"gopkg.in/yaml.v3"
)

type (
	BaseConfig[T BaseConfigData] struct {
		data      T
		fileName  string
		configDir string
	}
	BaseConfigData interface {
		Reset()
	}
)

func New[T BaseConfigData](defaultData T) (*BaseConfig[T], error) {
	configDir := getConfigDir()
	// Get filename from data type
	t := fmt.Sprintf("%T", defaultData)
	w := strings.Split(t, ".")
	if len(w) < 2 {
		panic(fmt.Sprintf("Invalid type for BaseConfig: %T", defaultData))
	}
	configDir = path.Join(configDir, "gdm")
	if err := pathtools.CreatePath(configDir); err != nil {
		return nil, err
	}
	fileName := path.Join(configDir, w[len(w)-1]+".yml")
	var data T = defaultData
	if content, err := os.ReadFile(fileName); err == nil {
		var readData T
		if err = yaml.Unmarshal(content, &readData); err == nil {
			data = readData
		}
	}
	data.Reset()
	return &BaseConfig[T]{
		data:      data,
		fileName:  fileName,
		configDir: configDir,
	}, nil
}

func (bc *BaseConfig[T]) Save() error {
	if content, err := yaml.Marshal(bc.data); err == nil {
		return os.WriteFile(bc.fileName, content, 0644)
	} else {
		return err
	}
}

func (bc *BaseConfig[T]) FileName() string {
	return bc.fileName
}

func (bc *BaseConfig[T]) ConfigDir() string {
	return bc.configDir
}

func (bc *BaseConfig[T]) Data() T {
	return bc.data
}
