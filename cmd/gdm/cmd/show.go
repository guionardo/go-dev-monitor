package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/guionardo/go-dev-monitor/internal/show"
	pathtools "github.com/guionardo/go/pkg/path_tools"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show current repository data",
	RunE:  runShow,
}

func runShow(command *cobra.Command, args []string) (err error) {
	currentDir := "."
	if len(args) > 0 {
		currentDir = args[0]
	}
	currentDir, err = filepath.Abs(currentDir)
	if err != nil {
		return fmt.Errorf("invalid path %s - %v", args[0], err)
	}

	if !pathtools.DirExists(currentDir) {
		return fmt.Errorf("path not found %s", currentDir)
	}
	return show.Show(currentDir)
}
