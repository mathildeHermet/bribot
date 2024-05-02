package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Sender struct {
	WebhookURL string
	Message    *Message
}

func (d *Sender) Send() error {
	var errMsg string
	jsonData, err := json.Marshal(d.Message)

	if err != nil {
		errMsg = fmt.Sprintf("Error encoding JSON: %v", err)
		return errors.New(errMsg)
	}

	req, err := http.NewRequest("POST", d.WebhookURL, bytes.NewBuffer(jsonData))

	if err != nil {
		errMsg = fmt.Sprintf("Error creating request: %v", err)
		return errors.New(errMsg)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		errMsg := fmt.Sprintf("Error sending request: %v", err)
		return errors.New(errMsg)
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	return nil
}
