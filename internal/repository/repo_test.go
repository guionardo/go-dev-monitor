package repository

import (
	"testing"

	"github.com/guionardo/go-dev-monitor/internal/utils/git"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("unexistent_folder", func(t *testing.T) {
		_, err := New("unexistent_folder", "localhost")
		assert.Error(t, err, "New() should return an error for unexistent folder")
	})
	t.Run("existent_folder_without_git", func(t *testing.T) {
		tmp := t.TempDir()
		_, err := New(tmp, "localhost")
		assert.Error(t, err, "New() should return an error for folder without git repository")
	})
	t.Run("existent_folder_with_git", func(t *testing.T) {
		root, err := git.GetGitRoot(".")
		assert.NoError(t, err, "Failed to get git root")
		repo, err := New(root, "localhost")
		assert.NoErrorf(t, err, "New() should not return an error for folder %s with git repository", root)
		assert.NotEmpty(t, repo.CurrentBranch)
		j, err := repo.MarshalJSON()
		// j, err := json.Marshal(repo)
		assert.NoError(t, err)
		t.Logf("JSON:\n%s", string(j))
		repo2 := Local{}
		err = repo2.UnmarshalJSON(j)
		assert.NoError(t, err)
		assert.Equal(t, repo.Origin, repo2.Origin)
	})
}
