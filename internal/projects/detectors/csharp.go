package detectors

import (
	"errors"
	"path"
)

const iconCSharp = "mdi mdi-language-csharp"

// DetectCSharp verify the files .csproj and .cs
func DetectCSharp(folderName string) (language string, projectName string, icon string, err error) {
	if firstFile := searchFiles(folderName, "*.csproj", "*.cs"); len(firstFile) > 0 {
		return "c#", path.Base(folderName), iconCSharp, nil
	}
	err = errors.New("not a C# project")
	return
}
