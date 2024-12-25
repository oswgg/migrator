package tools

import (
	"os"
	"strings"
)

func GetTxtValues(path string) (map[string]string, error) {
	var mappedValues = map[string]string{}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		var slicedLine = strings.SplitN(line, "=", 2)
		key := &slicedLine[0]
		value := &slicedLine[1]

		mappedValues[*key] = *value
	}

	return mappedValues, nil
}
