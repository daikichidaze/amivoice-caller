package common

import (
    "os"
)

func loadAPIKey(filePath string) (string, error) {
	apiKey, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(apiKey), nil
}
