package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Check if an audio file path is provided as an argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path to audio file>")
		return
	}

	audioFilePath := os.Args[1] // Get the path to the audio file from command line arguments

	// Read the API key from a local file named "APIKEY"
	apiKey, err := os.ReadFile("APIKEY") // Local file name containing the API key
	if err != nil {
		fmt.Println("Error reading API key file:", err)
		return
	}

	url := "https://acp-api-async.amivoice.com/v1/recognitions"
	method := "POST"

	// Open the audio file
	file, err := os.Open(audioFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add form fields
	_ = writer.WriteField("d", "-a-general")
	_ = writer.WriteField("u", string(apiKey)) // API key

	// Add the file part to the form
	part, err := writer.CreateFormFile("a", "test.wav")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the writer to finalize the multipart message
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create and send the request
	client := &http.Client{}
	req, err := http.NewRequest(method, url, &requestBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
