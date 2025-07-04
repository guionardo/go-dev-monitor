package cmd

import (
	"fmt"

	"github.com/guionardo/go-dev-monitor/internal"
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/debug"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gdm",
		Short:   "Go development monitor",
		Version: internal.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if enableDebug, err := cmd.Flags().GetBool("debug"); err == nil {
				debug.SetDebug(enableDebug)
			}
			if dataFolder, err := cmd.Flags().GetString("data"); err == nil && len(dataFolder) > 0 {
				config.SetConfigDir(dataFolder)
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debugging messages")
	defCfgDir, err := config.GetDefaultConfigDir()
	if err != nil {
		defCfgDir = fmt.Sprintf("* UNABLE TO DETECT USER CONFIG DIR: %v *", err)
	}
	rootCmd.PersistentFlags().String("data", "", fmt.Sprintf("data folder (default is %s)", defCfgDir))
	rootCmd.AddCommand(agentCmd, serverCmd, showCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
