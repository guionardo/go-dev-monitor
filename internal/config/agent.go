package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	pathtools "github.com/guionardo/go-dev-monitor/internal/utils/path_tools"
)

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

func (c *Agent) AddRoot(root string) error {
	if !pathtools.DirExists(root) {
		return fmt.Errorf("root not found %s", root)
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed getting absolute path of %s - %w", root, err)
	}
	if !slices.Contains(c.Roots, absRoot) {
		c.Roots = append(c.Roots, absRoot)
	}
	return nil
}

func (c *Agent) RemoveRoot(root string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed getting absolute path of %s - %w", root, err)
	}
	if p := slices.Index(c.Roots, absRoot); p >= 0 {
		c.Roots = append(c.Roots[0:p], c.Roots[p+1:]...)
	}
	return nil
}

func (c *Agent) SetServer(server string) error {
	url, err := url.Parse(server)
	if err != nil {
		return fmt.Errorf("failed setting server URI %s - %w", server, err)
	}
	if len(url.Scheme) == 0 {
		url.Scheme = "https"
	}
	if len(url.Host) == 0 {
		return fmt.Errorf("failed setting server URI %s - missing host", server)
	}
	c.ServerAddress = url.String()
	return nil
}
