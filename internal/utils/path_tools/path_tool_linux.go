package pathtools

import (
	"os"
	"syscall"
)

func CreatePath(path string) error {
	if DirExists(path) {
		return nil
	}
	oldmask := syscall.Umask(0)
	defer syscall.Umask(oldmask)
	return os.MkdirAll(path, os.ModeSticky|os.ModePerm)
}
