package finder

import (
	"path"
	"testing"

	utils_git "github.com/guionardo/go-dev-monitor/internal/utils/git"
	"github.com/stretchr/testify/assert"
)

func TestFindRepositories(t *testing.T) {
	devFolder, err := utils_git.GetGitRoot(".")
	if err != nil {
		t.Fatalf("Failed to get git root: %v", err)
	}
	devFolder = path.Dir(devFolder)
	t.Logf("Searching for projects in folder: %s with max depth 2", devFolder)

	folders := make([]string, 0, 10)
	for folder := range FindRepositories(devFolder, 2) {
		folders = append(folders, folder)
	}
	assert.Greater(t, len(folders), 0)
}
