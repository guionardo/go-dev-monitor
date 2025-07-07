package internal

import (
	_ "embed"
	"log/slog"
	"strings"

	"github.com/guionardo/go-dev-monitor/internal/logging"
)

//go:embed metadata.txt
var metadata string

var Version string = "dev"

func init() {
	Version = strings.TrimSpace(strings.ReplaceAll(metadata, "\n", ""))
	logging.Debug("Metadata version", slog.String("version", Version), slog.String("metadata", metadata))
}
