package show

import (
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/logging"
)

func Show(localFolder string) error {
	mainCfg := config.NewConfig()
	agentCfg := mainCfg.Agent
	if agentCfg == nil {
		agentCfg = config.NewAgentConfig()
		mainCfg.Agent = agentCfg
		if err := mainCfg.Save(); err != nil {
			logging.Error("saving config", err)
		}
	}
	agentCfg.Reset()

	return display(agentCfg, localFolder)
}
