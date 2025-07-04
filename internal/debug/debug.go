package debug

import (
	"log/slog"
	"os"
	"sync"
)

var (
	isDebug bool
	lock    sync.RWMutex
	logger  *slog.Logger
)

func init() {
	SetDebug(false)
}

func SetDebug(debug bool) {
	lock.Lock()
	isDebug = debug
	var level slog.Level
	if debug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	logger.Debug("Logging", slog.Bool("debug", debug))

	lock.Unlock()
}

func IsDebug() bool {
	lock.RLock()
	defer lock.RUnlock()
	return isDebug
}

func Log() *slog.Logger {
	lock.RLock()
	defer lock.RUnlock()
	return logger
}
