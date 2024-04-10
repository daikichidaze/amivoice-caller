package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

    "github.com/daikichidaze/amivoice-caller/pkg/common"
)

func fetchResults(sessionID, appKey string) ([]byte, error) {
	url := fmt.Sprintf("https://acp-api-async.amivoice.com/v1/recognitions/%s", sessionID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+appKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func formatAndSaveText(data []byte, outputPath string) error {
	// JSONデータをそのままファイルに保存
	return ioutil.WriteFile(outputPath, data, 0644)
}


func main() {
	// 仮定: "responce.json" からセッションIDを取得します。
	responseFilePath := "./response.json"

    apiKey, err := common.LoadAPIKey("./APIKEY")
	if err != nil {
		fmt.Println("Error loading API key:", err)
		return
	}

	responseData, err := ioutil.ReadFile(responseFilePath)
	if err != nil {
		fmt.Println("Error reading response file:", err)
		return
	}

	var response struct {
		SessionID string `json:"sessionid"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		fmt.Println("Error unmarshalling response data:", err)
		return
	}

	resultData, err := fetchResults(response.SessionID, apiKey)
	if err != nil {
		fmt.Println("Error fetching results:", err)
		return
	}

	if err := formatAndSaveText(resultData, "./result.json"); err != nil {
		fmt.Println("Error saving JSON data:", err)
		return
	}

	fmt.Println("JSON data saved successfully")
}
