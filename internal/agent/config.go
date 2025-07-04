package agent

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/guionardo/go-dev-monitor/internal/config"
)

type (
	ConfigData struct {
		Roots         []string `yaml:"roots"`
		MaxFolderDept int      `yaml:"max_folder_dept"`
		ServerAddress string   `yaml:"server_address"`
		Hostname      string   `yaml:"hostname"`
		configDir     string
	}
	AgentConfig config.BaseConfig[*ConfigData]
)

// type Config struct {
// 	Roots         []string `yaml:"roots"`
// 	MaxFolderDept int      `yaml:"max_folder_dept"`
// 	ServerAddress string   `yaml:"server_address"`
// 	Hostname      string   `yaml:"hostname"`
// 	configDir     string
// }

func (c *ConfigData) AddRoots(root string) {
	abs, _ := filepath.Abs(root)

	if !slices.Contains(c.Roots, abs) {
		c.Roots = append(c.Roots, abs)
	}
}

func (c *ConfigData) RemoveRoot(root string) {
	abs, _ := filepath.Abs(root)

	p := slices.Index(c.Roots, abs)
	if p == -1 {
		return
	}
	if len(c.Roots) == 1 {
		c.Roots = []string{}
		return
	}
	if p == 0 {
		c.Roots = c.Roots[1:]
	} else if p == len(c.Roots)-1 {
		c.Roots = c.Roots[0:p]
	} else {
		c.Roots = append(c.Roots[0:p], c.Roots[p+1:len(c.Roots)]...)
	}

}

func (c *ConfigData) Reset() {
	if c.Roots == nil {
		c.Roots = make([]string, 0)
	}
	if c.MaxFolderDept <= 0 {
		c.MaxFolderDept = 2
	}
	if len(c.ServerAddress) == 0 {
		c.ServerAddress = "http://localhost:3800"
	}
	if len(c.Hostname) == 0 {
		hostname, err := os.Hostname()
		if err == nil {
			c.Hostname = hostname
		} else {
			c.Hostname = "unknown"
		}
	}
}

func (c *ConfigData) GetRoots() []string {
	return c.Roots
}
