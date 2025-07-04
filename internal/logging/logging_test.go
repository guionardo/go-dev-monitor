package logging

import (
	"bytes"
	"io"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lock sync.Mutex

func fakeStdout() func() string {
	lock.Lock()
	r, w, _ := os.Pipe()
	stdOut = w
	return func() string {
		_ = w.Close()
		var capturedOutput bytes.Buffer
		_, _ = io.Copy(&capturedOutput, r)
		stdOut = os.Stdout
		lock.Unlock()
		return capturedOutput.String()
	}
}

func TestSetup(t *testing.T) {
	t.Run("default_stdout", func(t *testing.T) {
		read := fakeStdout()
		err := Setup(false, "")
		assert.NoError(t, err)
		logger.Info("default_stdout.1")
		logger.Debug("default_stdout.2")
		out := read()
		assert.Contains(t, out, "default_stdout.1")
		assert.NotContains(t, out, "default_stdout.2")
	})

	t.Run("default_stdout_debug", func(t *testing.T) {
		read := fakeStdout()
		err := Setup(true, "")
		assert.NoError(t, err)
		logger.Info("default_stdout_debug.1")
		logger.Debug("default_stdout_debug.2")
		out := read()
		assert.Contains(t, out, "default_stdout_debug.1")
		assert.Contains(t, out, "default_stdout_debug.2")
	})

	t.Run("default_logfile", func(t *testing.T) {
		logFile := path.Join(t.TempDir(), "log.log")
		read := fakeStdout()
		err := Setup(false, logFile)
		assert.NoError(t, err)
		logger.Info("default_stdout.1")
		logger.Debug("default_stdout.2")
		Close()
		outStd := read()
		out, err := os.ReadFile(currentlogFileName)
		assert.NoError(t, err)
		assert.Equal(t, outStd, string(out))
		assert.Contains(t, string(out), "default_stdout.1")
		assert.NotContains(t, string(out), "default_stdout.2")
	})

}
