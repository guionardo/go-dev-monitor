package show

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func Test_removeColors(t *testing.T) {
	color.NoColor = false

	text := fmt.Sprintf("Name: %s - Age: %s", color.RedString("%s", "Guionardo"), color.BlueString("%d", 48))
	assert.Equal(t, "Name: \x1b[31mGuionardo\x1b[0m - Age: \x1b[34m48\x1b[0m", text)

	noColorText := removeColors(text)
	assert.Equal(t, "Name: Guionardo - Age: 48", noColorText)

}
