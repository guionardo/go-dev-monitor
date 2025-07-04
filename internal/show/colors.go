package show

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type (
	Color uint8
)

const (
	White Color = iota
	Green
	Blue
	Yellow
	Cyan
	Red
	Grey
	undefinedColor
	link
)

var colorFuncs = map[Color]func(format string, args ...any) string{
	White:          color.WhiteString,
	Green:          color.GreenString,
	Blue:           color.BlueString,
	Yellow:         color.YellowString,
	Cyan:           color.CyanString,
	Red:            color.RedString,
	undefinedColor: fmt.Sprintf,
	// link:           consoleLink,
}

func consoleLink(format string, args ...any) string {
	//\e]8;;http://example.com\e\\This is a link\e]8;;\e\\\n
	// https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
	const startURI = "\x1b]8;;"
	const endURI = "\x1b"
	format = fmt.Sprintf("%s%s%s", startURI, format, endURI)
	return fmt.Sprintf(format, args...)
}

func (c Color) String() string {
	return fmt.Sprintf("{COLOR:%d}", c)
}

func colorPrefixSuffix(cf func(format string, args ...any) string) func() (prefix, suffix string) {
	return sync.OnceValues(func() (string, string) {
		s := cf("%s", "_COLOR_")
		ps := strings.SplitN(s, "_COLOR_", 2)
		return ps[0], ps[1]
	})
}
func removeColors(s string) string {
	for _, cf := range colorFuncs {
		prefix, suffix := colorPrefixSuffix(cf)()
		s = strings.ReplaceAll(s, prefix, "")
		s = strings.ReplaceAll(s, suffix, "")
	}
	return s
}
