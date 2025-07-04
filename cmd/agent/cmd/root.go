package cmd

import (
	"github.com/guionardo/go-dev-monitor/internal"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gdm",
		Short:   "",
		Version: internal.Version,
	}
)

func init() {
	rootCmd.AddCommand(runCmd, configCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
