package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
    "path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path to audio file>")
		return
	}
	audioFilePath := os.Args[1]

	apiKey, err := loadAPIKey("APIKEY")
	if err != nil {
		fmt.Println("Error loading API key:", err)
		return
	}

	requestBody, contentType, err := createMultiPartRequest(audioFilePath, apiKey)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	response, err := sendRequest("https://acp-api-async.amivoice.com/v1/recognitions", requestBody, contentType)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if err := saveResponseToFile(response, "response.txt"); err != nil {
		fmt.Println("Error saving response to file:", err)
		return
	}

	fmt.Println("Response saved to response.txt")
}

func loadAPIKey(filePath string) (string, error) {
	apiKey, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(apiKey), nil
}

func createMultiPartRequest(audioFilePath, apiKey string) (bytes.Buffer, string, error) {
	file, err := os.Open(audioFilePath)
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("d", "-a-general")
	_ = writer.WriteField("u", apiKey)

	// Extract the filename from the audioFilePath
	_, fileName := filepath.Split(audioFilePath)

	part, err := writer.CreateFormFile("a", fileName)
	if err != nil {
		return bytes.Buffer{}, "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return bytes.Buffer{}, "", err
	}

	err = writer.Close()
	if err != nil {
		return bytes.Buffer{}, "", err
	}

	return requestBody, writer.FormDataContentType(), nil
}


func sendRequest(url string, requestBody bytes.Buffer, contentType string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func saveResponseToFile(response, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(response)
	return err
}
