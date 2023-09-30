package pkg

import (
	"errors"
	"io/fs"
	"os"
)

// FileExist checks if a file exists at filePath.
func FileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return false
	}

	return true
}
