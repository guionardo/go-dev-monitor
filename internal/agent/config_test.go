package agent

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Roots(t *testing.T) {

	t.Run("add_roots", func(t *testing.T) {
		c := &ConfigData{}
		c.Reset()
		c.AddRoots("root_1")
		c.AddRoots("root_2")
		c.AddRoots("root_2")

		assert.Len(t, c.Roots, 2)
	})
	t.Run("del_roots_first", func(t *testing.T) {
		c := &ConfigData{}
		c.Reset()
		c.AddRoots("root_1")
		c.AddRoots("root_2")
		c.RemoveRoot("root_1")
		assert.Len(t, c.Roots, 1)
		assert.True(t, strings.HasSuffix(c.Roots[0], "root_2"))
	})

	t.Run("del_roots_last", func(t *testing.T) {
		c := &ConfigData{}
		c.Reset()
		c.AddRoots("root_1")
		c.AddRoots("root_2")
		c.RemoveRoot("root_2")
		assert.Len(t, c.Roots, 1)
		assert.True(t, strings.HasSuffix(c.Roots[0], "root_1"))
	})
	t.Run("del_roots_middle", func(t *testing.T) {
		c := &ConfigData{}
		c.Reset()
		c.AddRoots("root_1")
		c.AddRoots("root_2")
		c.AddRoots("root_3")
		c.AddRoots("root_4")
		c.RemoveRoot("root_2")
		assert.Len(t, c.Roots, 3)
		assert.True(t, strings.HasSuffix(c.Roots[0], "root_1"))
	})

}
