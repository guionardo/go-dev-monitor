package store

import (
	"testing"

	"github.com/guionardo/go-dev-monitor/internal/repository"
	"github.com/guionardo/go-dev-monitor/internal/utils/git"
	"github.com/stretchr/testify/assert"
)

func TestNewSqliteStore(t *testing.T) {
	tmp := t.TempDir()
	store, err := NewSqliteStore(tmp)
	if !assert.NoError(t, err) {
		return
	}

	root, _ := git.GetGitRoot(".")
	repo, err := repository.New(root, "localhost")
	if !assert.NoError(t, err) {
		return
	}

	err = store.Post("localhost", repo)
	if !assert.NoError(t, err) {
		return
	}

	summary, err := store.GetSummary()
	if !assert.NoError(t, err) {
		return
	}

	assert.Len(t, summary, 1)

}
