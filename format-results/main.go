package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Result struct {
	Label string `json:"label"`
	Text  string `json:"text"`
}

type Segment struct {
	Results []Result `json:"results"`
}

type Data struct {
	SessionID string    `json:"session_id"`
	Segments  []Segment `json:"segments"`
}

var nameDict = map[string]string{
	"speaker0": "吉岡",
	"speaker1": "竹屋さん",
}

func readJSONFile(filePath string) (Data, error) {
	var data Data
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return data, fmt.Errorf("error reading JSON file: %v", err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return data, fmt.Errorf("error parsing JSON: %v", err)
	}

	return data, nil
}

func formatSegments(segments []Segment, nameDict map[string]string) ([]string, error) {
	var textList []string
	lastSpoken := ""

	for _, seg := range segments {
		if len(seg.Results) > 1 {
			return nil, fmt.Errorf("error: long segment")
		}

		result := seg.Results[0]
		if result.Text == "" {
			continue
		}

		spoken := getSpeaker(result, nameDict)
		segTxt := formatSegmentText(spoken, lastSpoken, result.Text)
		lastSpoken = spoken

		textList = append(textList, segTxt)
	}

	return textList, nil
}

func getSpeaker(result Result, nameDict map[string]string) string {
	for _, token := range []Result{result} {
		if token.Label != "" {
			return nameDict[token.Label]
		}
	}
	return ""
}

func formatSegmentText(spoken, lastSpoken, text string) string {
	if lastSpoken != spoken {
		return fmt.Sprintf("\n%s: %s\n", spoken, text)
	}
	return fmt.Sprintf("%s\n", text)
}

func writeToFile(outputPath string, textList []string) error {
	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening output file: %v", err)
	}
	defer outputFile.Close()

	for _, item := range textList {
		_, err = outputFile.WriteString(item)
		if err != nil {
			return fmt.Errorf("error writing to output file: %v", err)
		}
	}

	return nil
}

func reverseNameDict(original map[string]string) map[string]string {
	reversed := make(map[string]string)
	for k, v := range original {
		if k == "speaker0" {
			reversed["speaker1"] = v
		} else if k == "speaker1" {
			reversed["speaker0"] = v
		}
	}
	return reversed
}

func main() {
	var suffix string
	flag.StringVar(&suffix, "suffix", "", "Suffix for the output file")
	reverseOrder := flag.Bool("reverse", false, "Reverse the order of nameDict")
	flag.Parse()

	filePath := "./../result.json"
	defaultOutputPath := "./formated-result.txt"
	outputPath := defaultOutputPath

	if suffix != "" {
		outputPath = fmt.Sprintf("./formated-result_%s.txt", suffix)
	}

	data, err := readJSONFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Start formatting results of %s\n", data.SessionID)

	selectedNameDict := nameDict
	if *reverseOrder {
		selectedNameDict = reverseNameDict(nameDict)
	}

	textList, err := formatSegments(data.Segments, selectedNameDict)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = writeToFile(outputPath, textList)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Finished formatting results of %s\n", data.SessionID)
}
