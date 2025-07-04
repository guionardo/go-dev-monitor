package cmd

import (
	"errors"
	"fmt"

	"github.com/guionardo/go-dev-monitor/internal/agent"
	"github.com/spf13/cobra"
)

var (
	agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "agent runner for cronjob",
		RunE:  runAgent,
	}
)

func init() {

	agentRunCmd := &cobra.Command{
		Use:   "run",
		Short: "fetch repositories from roots and publish to server",
		RunE:  runAgent,
	}
	agentAddRootCmd := &cobra.Command{
		Use:   "add",
		Short: "add folder to roots list",
		RunE:  agentAddRoot,
	}
	agentRemoveRootCmd := &cobra.Command{
		Use:   "del",
		Short: "remove folder from roots list",
		RunE:  agentDeleteRoot,
	}
	agentListRootsCmd := &cobra.Command{
		Use:   "list",
		Short: "list roots",
		RunE:  agentListRoots,
	}
	agentCmd.AddCommand(agentRunCmd, agentAddRootCmd, agentRemoveRootCmd, agentListRootsCmd)
}

func agentAddRoot(command *cobra.Command, args []string) error {
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

func agentDeleteRoot(command *cobra.Command, args []string) error {
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

func agentListRoots(command *cobra.Command, args []string) error {
	agent, err := agent.New()
	if err != nil {
		return err
	}

	for _, r := range agent.GetRoots() {
		fmt.Println(r)
	}
	return nil
}

func runAgent(command *cobra.Command, args []string) error {
	agent, err := agent.New()
	if err == nil {
		err = agent.Run()
	}
	return err
}
