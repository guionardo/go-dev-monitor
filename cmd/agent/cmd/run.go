package cmd

import (
	"github.com/guionardo/go-dev-monitor/internal/agent"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:  "run",
	RunE: run,
}

func run(command *cobra.Command, args []string) error {
	agent, err := agent.New()
	if err == nil {
		err = agent.Run()
	}
	return err
}
