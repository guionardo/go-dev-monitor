package detectors

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

const iconGolang = "fa-brands fa-golang"

// DetectGo verify the go.mod file
func DetectGo(folderName string) (language string, projectName string, icon string, err error) {
	content, err := os.ReadFile(path.Join(folderName, "go.mod"))
	if err != nil {
		err = fmt.Errorf("not a go project - %w", err)
		return
	}
	for line := range strings.SplitSeq(string(content), "\n") {
		// module github.com/guionardo/go-dev-monitor
		if w := strings.Split(line, " "); len(w) == 2 && w[0] == "module" {
			w = strings.Split(w[1], "/")
			return "golang", w[len(w)-1], iconGolang, nil
		}
	}
	err = errors.New("not a go project")
	return
}
