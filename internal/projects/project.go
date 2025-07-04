package projects

import (
	"fmt"
	"os"

	"github.com/guionardo/go-dev-monitor/internal/projects/detectors"
)

type (
	Project struct {
		FolderName  string `json:"folder_name"`
		Language    string `json:"language"`
		ProjectName string `json:"project_name"`
		Icon        string `json:"icon"`
	}
	DetectFunc func(folderName string) (language string, projectName string, icon string, err error)
)

var det = []DetectFunc{
	detectors.DetectGo,
	detectors.DetectPython,
	detectors.DetectCSharp,
}

func New(folderName string) (*Project, error) {
	if stat, err := os.Stat(folderName); err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("folder is invalid or not found: %s", folderName)
	}
	var language, projectName, icon string
	var err error
	for _, df := range det {
		if language, projectName, icon, err = df(folderName); err == nil {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if len(language) == 0 {
		language = "generic"
		icon = "fa-solid fa-file-lines"
	}

	return &Project{
		FolderName:  folderName,
		Language:    language,
		ProjectName: projectName,
		Icon:        icon,
	}, nil
}
