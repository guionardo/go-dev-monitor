package cmd

import (
	"fmt"

	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/server"
	pathtools "github.com/guionardo/go-dev-monitor/internal/utils/path_tools"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:  "run",
		RunE: runServer,
	}
)

func init() {
	
	serverCmd.Flags().Bool("debug", false, "enable debugging messages")
	serverCmd.Flags().String("data", "", "data folder (default is ~/.config/go_dev_monitor)")
}

func runServer(command *cobra.Command, _ []string) error {
	if dataFolder, _ := command.Flags().GetString("data"); len(dataFolder) > 0 {
		if !pathtools.DirExists(dataFolder) {
			return fmt.Errorf("data folder not found: %s", dataFolder)
		}
		config.SetConfigDir(dataFolder)
	}
	server, err := server.New()
	if err == nil {
		server.Run()
	}
	return err
}
