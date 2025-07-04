package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/guionardo/go-dev-monitor/internal/agent"
	"github.com/guionardo/go-dev-monitor/internal/config"
	pathtools "github.com/guionardo/go-dev-monitor/internal/utils/path_tools"
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

	agentShowCronCmd := &cobra.Command{
		Use:   "cron",
		Short: "show crontab usage",
		RunE:  agentShowCron,
	}
	agentInstallCmd := &cobra.Command{
		Use:   "install",
		Short: "install default production",
		RunE:  agentInstall,
	}
	agentInstallCmd.Flags().String("server", "https://devmon.guiosoft.info", "URL for go-dev-monitor server")

	defaultRoots := ""
	// Find dev folder in home directory
	if home, err := os.UserHomeDir(); err == nil {
		if pathtools.DirExists(path.Join(home, "dev")) {
			defaultRoots = path.Join(home, "dev")
		}
	}
	agentInstallCmd.Flags().String("root", defaultRoots, "root folders")

	hostname, _ := os.Hostname()
	agentInstallCmd.Flags().String("hostname", hostname, "machine host name")

	agentCmd.AddCommand(agentRunCmd, agentAddRootCmd, agentRemoveRootCmd, agentListRootsCmd, agentInstallCmd, agentShowCronCmd)
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

func agentInstall(command *cobra.Command, args []string) error {
	cfg := config.NewConfig()
	changed := false
	if root, err := command.Flags().GetString("root"); err == nil && len(root) > 0 {
		if err = cfg.Agent.AddRoot(root); err != nil {
			return err
		}
		changed = true
	}

	if server, err := command.Flags().GetString("server"); err == nil && len(server) > 0 {
		if err = cfg.Agent.SetServer(server); err != nil {
			return err
		}
		changed = true
	}

	if hostname, err := command.Flags().GetString("hostname"); err == nil && len(hostname) > 0 {
		cfg.Agent.Hostname = hostname
		changed = true
	}
	if changed {
		return cfg.Save()
	}

	return nil
}

func agentShowCron(command *cobra.Command, args []string) error {
	exec, err := os.Executable()
	if err != nil {
		return err
	}
	cfg := config.NewConfig()
	command.Printf(`Add the row above to your crontab file.
This will execute the collecting of all repositories every hour

> crontab -e

0 * * * * %s agent run --data "%s"

`, exec, cfg.GetConfigDir())
	return nil
}
