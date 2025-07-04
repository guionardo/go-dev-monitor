package detectors

import (
	"errors"
	"os"
	"path"
	"strings"
)

const iconPython = "fa-brands fa-python"

// DetectPython verify the pyprojet.toml or requirements.txt file
func DetectPython(folderName string) (language string, projectName string, icon string, err error) {
	if projectName, err = parsePyProject(folderName); err == nil {
		return "python", projectName, iconPython, nil
	}

	if projectName, err = parseRequirements(folderName); err == nil {
		return "python", projectName, iconPython, nil
	}

	err = errors.New("not a python project")
	return
}

func parsePyProject(folderName string) (projectName string, err error) {
	content, err := os.ReadFile(path.Join(folderName, "pyproject.toml"))
	if err != nil {
		return "", err
	}
	for line := range strings.SplitSeq(string(content), "\n") {
		line = strings.TrimSpace(line)
		if after, found := strings.CutPrefix(line, "name = "); found {
			return strings.TrimSpace(strings.ReplaceAll(after, "\"", "")), nil
		}
	}

	return "", errors.New("not a pyproject")

}

func parseRequirements(folderName string) (projectName string, err error) {
	if pythonFiles := searchFiles(folderName, "requirements.txt", "requirements*.txt", "*.py"); len(pythonFiles) > 0 {
		projectName = path.Base(folderName)
	}
	return
}
