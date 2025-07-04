package show

import (
	"fmt"
	"io"
	"iter"
	"os"
	"strings"

	"golang.org/x/term"
)

// const chars = "┌ └ ┐ ┘ ┼ ┬ ┴ ├ ┤ ─ │"

type (
	Panel struct {
		width   int
		Header  string
		Content string
		// c       *color.Color
	}
	Box struct {
		panels []*Panel
		width  int
	}
)

func NewBox() *Box {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = 80
	}
	return &Box{
		width: width,
	}
}

func (b *Box) AddPanel(format string, args ...any) *Panel {
	p := &Panel{
		width:  b.width - 4,
		Header: fmt.Sprintf(format, args...),
	}
	b.panels = append(b.panels, p)
	return p
}

func (b *Box) Write(w io.Writer) {
	if len(b.panels) == 0 {
		return
	}
	for i, p := range b.panels {

		switch i {
		case 0:
			p.ShowFirst(w)
		case len(b.panels) - 1:
			p.ShowLast(w)
		default:
			p.ShowMiddle(w)
		}
	}
	if len(b.panels) == 1 {
		padL("", "└", "┘", "─", b.width, w)
	}

}

func padL(s, prefix, suffix, char string, width int, w io.Writer) {
	l := len(removeColors(s))
	s = s[0:min(width-4, len(s))]

	if l > 0 {
		s = " " + s + " "
	}
	// pad := ""
	fmt.Fprintf(w, "%s%s\n", prefix, s)
	// if l < width-2 {
	// 	pad = strings.Repeat(char, width-2-l)
	// }
	// fmt.Fprintf(w, "%s%s%s%s\n", prefix, s, pad, suffix)
}

func (p *Panel) ShowFirst(w io.Writer) {

	padL(p.Header, "┌─", "┐", "─", p.width, w)
	p.showContent(w)
}

func (p *Panel) ShowMiddle(w io.Writer) {
	padL(p.Header, "├─", "┤", "─", p.width, w)
	p.showContent(w)
}

func (p *Panel) ShowLast(w io.Writer) {
	padL(p.Header, "├─", "┤", "─", p.width, w)
	p.showContent(w)
	padL("", "└─", "┘", "─", p.width, w)
}

func breakLine(line string, maxLength int) iter.Seq[string] {
	return func(yield func(string) bool) {
		var currentLine string
		for word := range strings.SplitSeq(line, " ") {
			if len(word)+len(currentLine) < maxLength {
				currentLine += " " + word
				continue
			}
			if !yield(currentLine) {
				return
			}
			currentLine = word
		}
		if len(currentLine) > 0 {
			yield(currentLine)
		}
	}
}

func (p *Panel) showContent(w io.Writer) {
	for line := range strings.SplitSeq(p.Content, "\n") {
		for subline := range breakLine(line, p.width) {
			padL(subline, "│  ", "│", " ", p.width, w)
		}
	}
}

func (p *Panel) AddContent(format string, args ...any) *Panel {
	var ln = ""
	if len(p.Content) > 0 {
		ln = "\n"
	}
	line := fmt.Sprintf(format, args...)
	p.Content = fmt.Sprintf("%s%s%s", p.Content, ln, line)
	return p
}
