package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAgentConfig(t *testing.T) {
	agent := NewAgentConfig()
	assert.NotNil(t, agent)

	assert.Error(t, agent.AddRoot("./unknow_test_root"))

	assert.NoError(t, agent.AddRoot("."))
	assert.Len(t, agent.Roots, 1)
	assert.NoError(t, agent.AddRoot(".."))
	assert.Len(t, agent.Roots, 2)

	assert.NoError(t, agent.RemoveRoot("."))
	assert.Len(t, agent.Roots, 1)
}
