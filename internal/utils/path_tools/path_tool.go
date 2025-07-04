package pathtools

import "os"

func DirExists(pathName string) bool {
	stat, err := os.Stat(pathName)
	return err == nil && stat.IsDir()
}
