package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"net/http"

	"github.com/guionardo/go-dev-monitor/internal/api"
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/repository"
	"github.com/guionardo/go-dev-monitor/internal/server"
	"github.com/guionardo/go-dev-monitor/internal/utils/finder"
)

type Agent struct {
	config *config.Config
}

func New() (*Agent, error) {
	cfg := config.NewConfig()

	return &Agent{
		config: cfg,
	}, nil
}

func (a *Agent) SaveConfigFile() error {
	return a.config.Save()
}

func (a *Agent) GetRoots() []string {
	return a.config.Agent.Roots
}

func (a *Agent) AddRoot(pathName string) error {
	return a.config.Agent.AddRoot(pathName)
}

func (a *Agent) RemoveRoot(pathName string) error {
	return a.config.Agent.RemoveRoot(pathName)
}

func (a *Agent) ProduceData() iter.Seq[*repository.Local] {
	return func(yield func(*repository.Local) bool) {
		for _, root := range a.config.Agent.Roots {
			for folder := range finder.FindRepositories(root, a.config.Agent.MaxFolderDept) {
				if repo, err := repository.New(folder, a.config.Agent.Hostname); err == nil {
					if !yield(repo) {
						return
					}
				}
			}
		}
	}
}

func (a *Agent) PublishData(data []*repository.Local) error {
	cfg := a.config.Agent
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
