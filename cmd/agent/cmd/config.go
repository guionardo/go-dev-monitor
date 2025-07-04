package cmd

import (
	"errors"

	"github.com/guionardo/go-dev-monitor/internal/agent"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use: "config",
}

func init() {

	cmdCfgAddRoot := &cobra.Command{
		Use:  "add-root",
		RunE: configAddRoot,
	}
	cmdCfgRemRoot := &cobra.Command{
		Use:  "del-root",
		RunE: configDelRoot,
	}
	configCmd.AddCommand(cmdCfgAddRoot, cmdCfgRemRoot)
}

func configAddRoot(command *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("expected folders for roots")
	}
	agent, err := agent.New()
	if err != nil {
		return err
	}
	for _, root := range args {
		if err := agent.AddRoot(root); err != nil {
			return err
		}
		command.Printf("Added root: %s\n", root)
	}
	return agent.SaveConfigFile()
}

func configDelRoot(command *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("expected folders for roots")
	}
	agent, err := agent.New()
	if err != nil {
		return err
	}
	for _, root := range args {
		if err := agent.RemoveRoot(root); err != nil {
			return err
		}
	}
	return agent.SaveConfigFile()
}
