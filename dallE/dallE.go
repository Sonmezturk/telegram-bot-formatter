package dallE

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type DallePayload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type Response struct {
	Created int `json:"created"`
	Data    []struct {
		RevisedPrompt string `json:"revised_prompt"`
		URL           string `json:"url"`
	} `json:"data"`
}

// GenerateImage calls the OpenAI API to generate an image
func GenerateImage(prompt string, n int, size string) (*Response, error) {

	url := "https://api.openai.com/v1/images/generations"
	method := "POST"

	// Create an instance of the payload with your details
	payload := DallePayload{
		Model:  "dall-e-3",
		Prompt: prompt,
		N:      1,
		Size:   "1024x1024",
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a reader from the JSON bytes
	payloadReader := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	bearerToken := fmt.Sprintf(`Bearer %s`, os.Getenv("OPENAI_KEY"))
	req.Header.Add("Authorization", bearerToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return &response, nil
}
