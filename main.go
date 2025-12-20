package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

type SignRequest struct {
	Data string `json:"data"`
}

type SignResponse struct {
	Signature string `json:"signature"`
}

func signData() error {
	// Create HTTP client with Unix socket transport
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/thunderwind/agent.sock")
			},
		},
		Timeout: 10 * time.Second,
	}

	// Prepare request body
	reqBody := SignRequest{
		Data: "SGVsbG8sIHdvcmxkIQ==", // base64 of "Hello, world!"
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make POST request
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"http://localhost/sign",
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var signResp SignResponse
	if err := json.NewDecoder(resp.Body).Decode(&signResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Save signature to file
	if err := os.WriteFile("/tmp/signature.txt", []byte(signResp.Signature), 0644); err != nil {
		return fmt.Errorf("failed to write signature file: %w", err)
	}

	fmt.Println("Successfully signed data and saved signature to /tmp/signature.txt")
	return nil
}

func main() {
	fmt.Println("HELLO from Thunderwind Enclave!")

	// Attempt to sign data via Unix socket
	if err := signData(); err != nil {
		fmt.Printf("Warning: signing failed: %v\n", err)
	}

	fmt.Println("Waiting 5 seconds...")
	time.Sleep(5 * time.Second)
	fmt.Println("Goodbye!")
}
