package tools

import "os"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func WriteFile(filePath, content string, perm os.FileMode) error {
	return os.WriteFile(filePath, []byte(content), perm)
}
