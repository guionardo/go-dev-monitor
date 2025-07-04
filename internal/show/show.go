package show

import (
	"github.com/guionardo/go-dev-monitor/internal/config"
)

func Show(localFolder string) error {
	mainCfg := config.NewConfig()
	agentCfg := mainCfg.Agent
	if agentCfg == nil {
		agentCfg = config.NewAgentConfig()
		mainCfg.Agent = agentCfg
		mainCfg.Save()
	}
	agentCfg.Reset()

	return display(agentCfg, localFolder)
}
