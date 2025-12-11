package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
)

// kestraExecuteResponse is a minimal shape of what Kestra returns on execution
type kestraExecuteResponse struct {
	ExecutionID string `json:"id"`
}

// triggerKestraWorkflow starts the Kestra workflow that runs Syft + Grype + AI.
func triggerKestraWorkflow(repoURL string) (string, error) {
	// Base URL where Kestra is running (inside Docker network)
	kestraBaseURL := os.Getenv("KESTRA_BASE_URL")
	if kestraBaseURL == "" {
		// Fallback if running outside Docker; adjust if needed
		kestraBaseURL = "http://localhost:8081"
	}

	kestraNamespace := os.Getenv("KESTRA_NAMESPACE")
	if kestraNamespace == "" {
		kestraNamespace = "autosec"
	}

	kestraFlowID := os.Getenv("KESTRA_WORKFLOW_ID")
	if kestraFlowID == "" {
		kestraFlowID = "autosec-scan"
	}

	// Correct Kestra executions API:
	// POST /api/v1/main/executions/{namespace}/{flowId}
	url := fmt.Sprintf("%s/api/v1/main/executions/%s/%s", kestraBaseURL, kestraNamespace, kestraFlowID)

	// Build multipart/form-data body with inputs (repo_url)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	field, err := writer.CreateFormField("repo_url")
	if err != nil {
		return "", fmt.Errorf("create form field: %w", err)
	}
	if _, err := field.Write([]byte(repoURL)); err != nil {
		return "", fmt.Errorf("write form field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// --- Authentication section ---

	// 1) Prefer API token if present
	if token := os.Getenv("KESTRA_API_TOKEN"); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else if user := os.Getenv("KESTRA_BASIC_USERNAME"); user != "" {
		// 2) Fall back to Basic Auth with username/password
		pass := os.Getenv("KESTRA_BASIC_PASSWORD")
		credential := user + ":" + pass
		encoded := base64.StdEncoding.EncodeToString([]byte(credential))
		req.Header.Set("Authorization", "Basic "+encoded)
	}

	// --- End authentication section ---

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call kestra: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("kestra returned status %d", resp.StatusCode)
	}

	var kestraResp kestraExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&kestraResp); err != nil {
		return "", fmt.Errorf("decode kestra response: %w", err)
	}

	if kestraResp.ExecutionID == "" {
		kestraResp.ExecutionID = "unknown-execution-id"
	}

	return kestraResp.ExecutionID, nil
}
