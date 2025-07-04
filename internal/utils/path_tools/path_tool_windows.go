package pathtools

import (
	"os"
)

func CreatePath(path string) error {
	if DirExists(path) {
		return nil
	}
	return os.Mkdir(path, os.ModeSticky|os.ModePerm)
}
