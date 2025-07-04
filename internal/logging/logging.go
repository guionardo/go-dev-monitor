package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	logWriter          io.Writer
	logFile            io.WriteCloser
	logger             *slog.Logger
	currentlogFileName string
	stdOut             io.Writer = os.Stdout
	setup                        = sync.OnceFunc(func() {
		if logger == nil {
			_ = Setup(true, "")
		}
	})
)

const head = "logging.setup"

func Setup(debug bool, logFileName string) (err error) {
	if len(logFileName) > 0 {
		lf, realLf, errFile := isValidLogFile(logFileName)
		if errFile != nil {
			err = errFile
			slog.Error(head, slog.Any("error", err))

		} else {
			logFile = lf
			logFileName = realLf
			logWriter = io.MultiWriter(stdOut, lf)
		}
	} else {
		logWriter = stdOut
	}
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(logWriter, &slog.HandlerOptions{Level: level}))
	logger.Debug(head, slog.Bool("debug", debug), slog.String("filename", logFileName))
	currentlogFileName = logFileName
	return nil
}

func SetupGin(router *gin.Engine) {
	router.Use(gin.LoggerWithWriter(logWriter))
}

func isValidLogFile(logFile string) (f io.WriteCloser, realLogFile string, err error) {
	lf, err := filepath.Abs(logFile)
	if err == nil {
		lf = getRotationLogFile(lf, 7)
		f, err = os.OpenFile(lf, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err == nil {
			_, err = f.Write([]byte{})
		}
	}
	if err != nil {
		if f != nil {
			_ = f.Close()
		}
		return nil, lf, fmt.Errorf("invalid log file %s - %w", logFile, err)
	}
	return f, lf, nil

}

func getRotationLogFile(logFile string, keepCount int) string {
	ext := filepath.Ext(logFile)
	logFilePrefix, _ := strings.CutSuffix(logFile, ext)
	filter := logFilePrefix + "*" + ext
	currentLogFile := logFilePrefix + "." + time.Now().Format("2006-05-04") + ext
	logFiles, err := filepath.Glob(filter)
	if err != nil {
		slog.Error(head, slog.Any("error", err))
		return currentLogFile
	}
	if len(logFiles) > keepCount {
		sort.Strings(logFiles)
		removed := logFiles[0:keepCount]
		for _, r := range removed {
			_ = os.Remove(r)
		}
		slog.Info(head, slog.Any("removed", err))
	}
	return currentLogFile
}

func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

func IsDebug() bool {
	setup()
	if logger != nil {
		return logger.Enabled(context.Background(), slog.LevelDebug)
	}
	return false
}

func With(args ...any) *slog.Logger {
	setup()
	return logger.With(args...)
}

func Info(msg string, args ...any) {
	setup()
	logger.Info(msg, args...)
}

func Error(msg string, err error, args ...any) {
	setup()
	logger.Error(msg, append(args, slog.Any("error", err))...)
}

func Warn(msg string, args ...any) {
	setup()
	logger.Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	setup()
	logger.Debug(msg, args...)
}
