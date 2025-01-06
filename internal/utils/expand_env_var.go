package utils

import (
	"errors"
	"os"
)

func ExpandEnvVar(s string) (string, error) {
	var look = func(s string) string {
		return os.Getenv(s)
	}
	value := os.Expand(s, look)
	if value == "" {
		return "", errors.New(s + " is required")
	}

	return value, nil
}
