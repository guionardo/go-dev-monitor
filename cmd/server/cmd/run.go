package cmd

import (
	"github.com/guionardo/go-dev-monitor/internal/server"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:  "run",
		RunE: run,
	}
)

func init() {
	runCmd.Flags().Bool("debug", false, "enable debugging messages")
}

func run(command *cobra.Command, args []string) error {

	server, err := server.New()
	if err == nil {
		server.Run()
	}
	return err
}
