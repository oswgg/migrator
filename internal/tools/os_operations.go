package tools

import "os"

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func WriteFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}
