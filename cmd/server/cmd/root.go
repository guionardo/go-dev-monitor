package cmd

import (
	"github.com/guionardo/go-dev-monitor/internal"
	"github.com/guionardo/go-dev-monitor/internal/debug"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gdm",
		Short:   "",
		Version: internal.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if enableDebug, err := cmd.Flags().GetBool("debug"); err == nil {
				debug.SetDebug(enableDebug)
			}
		},
	}
)

func init() {

	rootCmd.AddCommand(runCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
