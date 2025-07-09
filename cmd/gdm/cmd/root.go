package cmd

import (
	"fmt"
	"log/slog"

	"github.com/guionardo/go-dev-monitor/internal"
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/logging"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gdm",
		Short:   "Go development monitor",
		Version: internal.Version,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logFile, _ := cmd.Flags().GetString("log")
			enableDebug, _ := cmd.Flags().GetBool("debug")
			if err := logging.Setup(enableDebug, logFile); err != nil {
				return err
			}

			if dataFolder, err := cmd.Flags().GetString("data"); err == nil && len(dataFolder) > 0 {
				config.SetConfigDir(dataFolder)
			}
			logging.Debug("Metadata version", slog.String("version", internal.Version))
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().Bool("debug", false, "enable debugging messages")
	rootCmd.PersistentFlags().String("log", "", "log file")
	defCfgDir, err := config.GetDefaultConfigDir()
	if err != nil {
		defCfgDir = fmt.Sprintf("* UNABLE TO DETECT USER CONFIG DIR: %v *", err)
	}
	rootCmd.PersistentFlags().String("data", "", fmt.Sprintf("data folder (default is %s)", defCfgDir))
	rootCmd.AddCommand(agentCmd, serverCmd, showCmd)
}

func Execute() error {
	defer logging.Close()
	return rootCmd.Execute()
}
