package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type ConfigSample struct {
	Name    string
	Enabled bool
	When    time.Time
}

func (c *ConfigSample) Reset() {}

func TestNew(t *testing.T) {

	SetConfigDir(t.TempDir())
	cfg, err := New(&ConfigSample{})
	assert.NoError(t, err)
	assert.Equal(t, "", cfg.data.Name)

}
