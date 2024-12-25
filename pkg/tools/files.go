package tools

import (
	"fmt"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func ReadFile(path string) ([]byte, error) {
	exists := FileExists(path)
	if !exists {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}

	return os.ReadFile(path)
}

func WriteFile(filePath, content string, perm os.FileMode) error {
	return os.WriteFile(filePath, []byte(content), perm)
}
