package detectors

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func touch(t *testing.T, folder, fileName string) {
	assert.NoError(t, os.WriteFile(path.Join(folder, fileName), []byte{}, 0655))
}
func Test_searchFiles(t *testing.T) {
	tmp := t.TempDir()
	touch(t, tmp, "t1.py")
	touch(t, tmp, "t2.py")

	f := searchFiles(tmp, "t1.py")

	assert.True(t, strings.HasSuffix(f, "t1.py"))
	assert.True(t, strings.HasSuffix(searchFiles(tmp, "*.py"), "t1.py"))

}
