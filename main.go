package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

// Command-line arguments
var (
	filePath string
	token    string
)

func init() {
	flag.StringVar(&filePath, "file", "", "Path to the file")
	flag.StringVar(&token, "token", "", "OpenAI API token")
	flag.Parse()

	if filePath == "" || token == "" {
		color.Red("[!] Failed to start.")
		flag.Usage()
		os.Exit(1)
	}
}

// Step 2: Upload the training data file
func uploadFile() string {
	url := "https://api.openai.com/v1/files"
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		color.Red("Failed to create form file: %v", err)
		return ""
	}
	file, err := os.Open(filePath)
	if err != nil {
		color.Red("Failed to open file: %v", err)
		return ""
	}
	defer file.Close()
	_, err = io.Copy(part, file)
	if err != nil {
		color.Red("Failed to copy file data: %v", err)
		return ""
	}
	err = writer.Close()
	if err != nil {
		color.Red("Failed to close multipart writer: %v", err)
		return ""
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		color.Red("Failed to create HTTP request: %v", err)
		return ""
	}
	req.Header = headers
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Failed to send HTTP request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Failed to read response body: %v", err)
		return ""
	}

	fileID := gjson.GetBytes(responseBody, "id").String()
	return fileID
}

// Check if the uploaded file is processed
func isFileProcessed(fileID string) bool {
	url := fmt.Sprintf("https://api.openai.com/v1/files/%s", fileID)
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("Failed to create HTTP request: %v", err)
		return false
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Failed to send HTTP request: %v", err)
		return false
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Failed to read response body: %v", err)
		return false
	}

	return gjson.GetBytes(responseBody, "status").String() == "processed"
}

// Step 3: Create a fine-tuning job
func createFineTuningJob(fileID string) string {
	url := "https://api.openai.com/v1/fine_tuning/jobs"
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
		"Content-Type":  []string{"application/json"},
	}
	body := map[string]interface{}{
		"training_file": fileID,
		"model":         "gpt-3.5-turbo-0613",
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		color.Red("Failed to marshal JSON body: %v", err)
		return ""
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		color.Red("Failed to create HTTP request: %v", err)
		return ""
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Failed to send HTTP request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Failed to read response body: %v", err)
		return ""
	}

	jobID := gjson.GetBytes(responseBody, "id").String()
	return jobID
}

// Check if the fine-tuning job is ready
func isJobReady(jobID string) bool {
	url := fmt.Sprintf("https://api.openai.com/v1/fine_tuning/jobs/%s", jobID)
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("Failed to create HTTP request: %v", err)
		return false
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Failed to send HTTP request: %v", err)
		return false
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Failed to read response body: %v", err)
		return false
	}

	return gjson.GetBytes(responseBody, "status").String() == "succeeded"
}

// Wait until the fine-tuning job is ready
func waitForJob(jobID string) {
	startTime := time.Now() // Record the start time

	for !isJobReady(jobID) {
		elapsedTime := time.Since(startTime) // Calculate the elapsed time
		fmt.Printf("\r%s Waiting for fine-tuning job to complete... Elapsed Time: %.2f seconds", color.YellowString("[!]"), elapsedTime.Seconds())
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}

func main() {
	startTime := time.Now() // Record the start time
	fmt.Printf("%s Starting... Elapsed Time: 0.00 seconds\n", color.YellowString("[!]"))

	// Upload the training data file
	fileID := uploadFile()
	if fileID == "" {
		color.Red("File upload failed.")
		return
	}

	elapsedTime := time.Since(startTime) // Calculate the elapsed time
	fmt.Printf("%s File uploaded successfully. File ID: %s Elapsed Time: %.2f seconds\n", color.GreenString("[+]"), fileID, elapsedTime.Seconds())

	// Check if the uploaded file is processed
	startTime = time.Now() // Record the start time
	for !isFileProcessed(fileID) {
		elapsedTime := time.Since(startTime) // Calculate the elapsed time
		fmt.Printf("\r%s Waiting for file processing... Elapsed Time: %.2f seconds", color.YellowString("[!]"), elapsedTime.Seconds())
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
	fmt.Printf("%s File processing completed.\n", color.GreenString("[+]"))

	elapsedTime = time.Since(startTime) // Calculate the elapsed time
	fmt.Printf("%s Starting Fine-tuning Job... Elapsed Time: %.2f seconds\n", color.YellowString("[!]"), elapsedTime.Seconds())

	// Create the fine-tuning job
	jobID := createFineTuningJob(fileID)
	if jobID == "" {
		color.Red("Fine-tuning job creation failed.")
		return
	}

	elapsedTime = time.Since(startTime) // Calculate the elapsed time
	fmt.Printf("%s Fine-tuning job created successfully. Job ID: %s\n", color.GreenString("[+]"), jobID)

	// Wait for the fine-tuning job to complete
	startTime = time.Now() // Record the start time
	waitForJob(jobID)

	fmt.Printf("%s Fine-tuning job completed.\n", color.GreenString("[+]"))

	elapsedTime = time.Since(startTime) // Calculate the elapsed time
	fmt.Printf("%s Elapsed Time: %.2f seconds\n", color.YellowString("[!]"), elapsedTime.Seconds())

	// Get the fine-tuned model from the completed job
	url := fmt.Sprintf("https://api.openai.com/v1/fine_tuning/jobs/%s", jobID)
	headers := http.Header{
		"Authorization": []string{"Bearer " + token},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		color.Red("Failed to create HTTP request: %v", err)
		return
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		color.Red("Failed to send HTTP request: %v", err)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Failed to read response body: %v", err)
		return
	}

	fineTunedModel := gjson.GetBytes(responseBody, "fine_tuned_model").String()

	// Print the fine-tuned model
	if fineTunedModel != "" {
		fmt.Printf("\nFine-tuned model ready to use: %s\n", fineTunedModel)
	} else {
		fmt.Println("No fine-tuned model available.")
	}
}
