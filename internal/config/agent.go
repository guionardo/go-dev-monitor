package config

import "os"

type Agent struct {
	Roots         []string `yaml:"roots"`
	MaxFolderDept int      `yaml:"max_folder_dept"`
	ServerAddress string   `yaml:"server_address"`
	Hostname      string   `yaml:"hostname"`
}

func NewAgentConfig() *Agent {
	return (&Agent{}).Reset()
}

func (c *Agent) Reset() *Agent {
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
	return c
}
