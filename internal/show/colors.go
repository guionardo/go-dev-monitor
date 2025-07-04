package show

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type (
	Color       uint8
	ColorString struct {
		format      string
		colorFormat string
		args        []argColor
	}
	argColor struct {
		arg    any
		color  Color
		format string
	}
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
	link:           consoleLink,
}

func consoleLink(format string, args ...any) string {
	//\e]8;;http://example.com\e\\This is a link\e]8;;\e\\\n
	// https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
	const startURI = "\x1b]8;;"
	const endURI = "\x1b"
	format = fmt.Sprintf("%s%s%s", startURI, format, endURI)
	return fmt.Sprintf(format, args...)
}

func (ac argColor) String() string {
	return colorFuncs[ac.color](ac.format, ac.arg)
}

func (ac argColor) Raw() string {
	return fmt.Sprintf(ac.format, ac.arg)
}

func (c Color) String() string {
	return fmt.Sprintf("{COLOR:%d}", c)
}

func NewCS(format string, args ...any) *ColorString {
	// Extract all format declarations
	formats := make([]string, 0, len(args))
	colorFormat := format
	tmpFmt := format
	nextFmtIndex := -1
	for i, c := range tmpFmt {
		if c == '%' {
			nextFmtIndex = i + 1
			continue
		}
		if nextFmtIndex > 0 {
			formats = append(formats, tmpFmt[nextFmtIndex-1:nextFmtIndex+1])
			colorFormat = colorFormat[0:nextFmtIndex] + "s" + colorFormat[nextFmtIndex+1:]
			nextFmtIndex = -1
		}
	}
	for len(tmpFmt) > 0 {
		if p := strings.Index(tmpFmt, "%"); p >= 0 && p < len(tmpFmt)-1 {
			formats = append(formats, tmpFmt[p:p+2])
			tmpFmt = tmpFmt[p+1:]
		} else {
			tmpFmt = ""
		}
	}
	var (
		nArgs     = make([]argColor, 0, len(args))
		lastColor = undefinedColor
		lastArg   any
	)

	for _, arg := range args {
		if color, ok := arg.(Color); ok {
			if lastColor == undefinedColor {
				lastColor = color
			}
		} else {
			if lastArg == nil {
				lastArg = arg
			}
		}
		if lastColor != undefinedColor && lastArg != nil {
			nArgs = append(nArgs, argColor{
				arg:    lastArg,
				color:  lastColor,
				format: formats[len(nArgs)],
			})
			lastArg = nil
			lastColor = undefinedColor
		}
	}
	if lastArg != nil {
		nArgs = append(nArgs, argColor{
			arg:   lastArg,
			color: undefinedColor,
		})
	}
	cs := &ColorString{
		format:      format,
		colorFormat: colorFormat,
		args:        nArgs,
	}

	return cs
}

func (cs *ColorString) String() string {
	var args = make([]any, len(cs.args))
	for i, arg := range cs.args {
		args[i] = arg.String()
	}
	return fmt.Sprintf(cs.colorFormat, args...)
}

func (cs *ColorString) NoColor() string {
	var args = make([]any, len(cs.args))
	for i, arg := range cs.args {
		args[i] = arg.arg
	}
	return fmt.Sprintf(cs.format, args...)
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
