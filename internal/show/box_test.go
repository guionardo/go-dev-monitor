package show

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBox(t *testing.T) {
	box := NewBox()
	assert.NotNil(t, box)
	box.width = 30
	box.AddPanel("Box test").AddContent("O rato roeu a roupa do rei de Roma. O rei, de raiva, roeu o resto")
	box.AddPanel("Second panel").AddContent("data = 1234\nname=\"test\"")
	box.AddPanel("Last panel").AddContent("1234")
	w := bytes.NewBufferString("")
	box.Write(w)
	fmt.Println(w.String())

}

func Test_breakLine(t *testing.T) {
	longLine := "O rato roeu a roupa do rei de Roma. O rei, de raiva, roeu o resto"
	var lines []string
	for line := range breakLine(longLine, 12) {
		lines = append(lines, line)
	}

	assert.Equal(t, []string{
		"O rato roeu",
		"a roupa do",
		"rei de Roma.",
		"O rei, de",
		"raiva, roeu",
		"o resto",
	}, lines)

}

func Test_padL(t *testing.T) {
	fw := func(s, prefix, suffix, char string, n int) string {
		b := bytes.NewBufferString("")
		padL(s, prefix, suffix, char, n, b)
		return b.String()
	}
	assert.Equal(t, "| abcd. \n", fw("abcd.", "|", "|", " ", 10))
	assert.Equal(t, "| abc de \n", fw("abc def ghi", "|", "|", " ", 10))

}
