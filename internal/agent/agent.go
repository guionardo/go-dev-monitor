package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"net/http"
	"os"

	"github.com/guionardo/go-dev-monitor/internal/api"
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/repository"
	"github.com/guionardo/go-dev-monitor/internal/server"
	"github.com/guionardo/go-dev-monitor/internal/utils/finder"
)

type Agent struct {
	config *config.BaseConfig[*ConfigData]
}

func New() (*Agent, error) {
	cfg, err := config.New(&ConfigData{})
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: cfg,
	}, nil
}

func (a *Agent) SaveConfigFile() error {
	return a.config.Save()
}

func (a *Agent) GetRoots() []string {
	return a.config.Data().Roots
}

func (a *Agent) AddRoot(pathName string) error {
	if stat, err := os.Stat(pathName); err == nil && stat.IsDir() {
		a.config.Data().AddRoots(pathName)
		return nil
	}
	return fmt.Errorf("invalid root %s", pathName)
}

func (a *Agent) RemoveRoot(pathName string) error {
	a.config.Data().RemoveRoot(pathName)
	return nil
}

func (a *Agent) ProduceData() iter.Seq[*repository.Local] {
	return func(yield func(*repository.Local) bool) {
		for _, root := range a.config.Data().GetRoots() {
			for folder := range finder.FindRepositories(root, a.config.Data().MaxFolderDept) {
				if repo, err := repository.New(folder, a.config.Data().Hostname); err == nil {
					if !yield(repo) {
						return
					}
				}
			}
		}
	}
}

func (a *Agent) PublishData(data []*repository.Local) error {
	cfg := a.config.Data()
	requestData := api.AgentRequest{
		Hostname:     cfg.Hostname,
		Repositories: data,
	}
	body, err := json.Marshal(requestData)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/data", cfg.ServerAddress), reader)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(server.AppDataHeader, server.AppDataValue)
	req.Header.Add(server.HostNameHeader, cfg.Hostname)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected API response - %s", resp.Status)
	}
	return nil
}

func (a *Agent) Run() error {
	data := make([]*repository.Local, 0)
	for r := range a.ProduceData() {
		data = append(data, r)
	}
	if len(data) == 0 {
		return errors.New("no repositories found")
	}
	return a.PublishData(data)
}
