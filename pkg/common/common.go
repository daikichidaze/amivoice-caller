package common

import (
    "os"
)

func LoadAPIKey(filePath string) (string, error) {
	apiKey, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(apiKey), nil
}
