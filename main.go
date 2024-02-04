package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
    "strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path to audio file>")
		return
	}
	audioFilePath := os.Args[1]

	// Convert the audio file to MP3 format
	mp3FilePath, err := convertToMP3(audioFilePath)
	if err != nil {
		fmt.Println("Error converting file to MP3:", err)
		return
	}

	apiKey, err := loadAPIKey("APIKEY")
	if err != nil {
		fmt.Println("Error loading API key:", err)
		return
	}

	requestBody, contentType, err := createMultiPartRequest(mp3FilePath, apiKey)
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

// Other functions (loadAPIKey, createMultiPartRequest, sendRequest, saveResponseToFile) remain unchanged.

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

func convertToMP3(inputFilePath string) (string, error) {
	// Extract the file name without extension
	baseName := strings.TrimSuffix(inputFilePath, filepath.Ext(inputFilePath))
	outputFilePath := baseName + ".mp3"

	// Create the ffmpeg command to convert the input file to MP3 format
	cmd := exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-ar", "44100", "-ac", "2", "-b:a", "192k", outputFilePath)
	
	// Connect ffmpeg's standard output and standard error output to the current process's output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Run the ffmpeg command and return any errors
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Return the path of the converted MP3 file
	return outputFilePath, nil
}