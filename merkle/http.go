package merkle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func MakePost(url string, apiKey string, body interface{}, resp interface{}) error {
	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("error marshalling bundle: %v", err)
	}

	if body == nil {
		bodyBytes = []byte{}
	}

	buffer := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest("POST", url, buffer)

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	// read the whole body
	defer res.Body.Close()

	bodyRead, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	if res.StatusCode > 400 {
		return fmt.Errorf("error sending request: url=%s code=%s, body=%s", url, res.Status, bodyRead)
	}

	if resp == nil {
		return nil
	}

	err = json.Unmarshal(bodyRead, &resp)

	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	return nil
}

func MakeDel(url string, apiKey string, body interface{}, resp interface{}) error {
	bodyBytes, err := json.Marshal(body)

	if err != nil {
		return fmt.Errorf("error marshalling bundle: %v", err)
	}

	if body == nil {
		bodyBytes = []byte{}
	}

	buffer := bytes.NewBuffer(bodyBytes)

	req, err := http.NewRequest("DELETE", url, buffer)

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	// read the whole body
	defer res.Body.Close()

	bodyRead, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	if res.StatusCode > 400 {
		return fmt.Errorf("error sending request: url=%s, code=%s, body=%s", url, res.Status, bodyRead)
	}

	if resp == nil {
		return nil
	}

	err = json.Unmarshal(bodyRead, &resp)

	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	return nil
}

func MakeGet(url string, apiKey string, resp interface{}) error {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+apiKey)

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	if res.StatusCode > 400 {
		return fmt.Errorf("error sending request: code=%s", res.Status)
	}

	if resp == nil {
		return nil
	}

	err = json.NewDecoder(res.Body).Decode(&resp)

	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	return nil
}
