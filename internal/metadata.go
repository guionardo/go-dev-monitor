package internal

import (
	_ "embed"
	"strings"
)

//go:embed metadata.txt
var metadata string

var Version string = "dev"

func init() {
	Version = strings.TrimSpace(strings.ReplaceAll(metadata, "\n", ""))

}
