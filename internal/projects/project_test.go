package projects

import (
	"testing"

	"github.com/guionardo/go-dev-monitor/internal/utils/git"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	root, err := git.GetGitRoot(".")
	assert.NoError(t, err)
	p, err := New(root)
	assert.NoError(t, err)
	assert.Equal(t, "golang", p.Language)

	p, err = New("/home/guionardo/dev/planetapeia-desfiles")
	assert.NoError(t, err)
	assert.Equal(t, "python", p.Language)
}
