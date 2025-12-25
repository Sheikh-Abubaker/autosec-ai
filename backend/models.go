package main

import "time"

// ScanRequest is the payload from the frontend when user clicks "Scan & Auto-Fix"
type ScanRequest struct {
	RepoURL string `json:"repo_url" binding:"required"`
}

// ScanResponse is a simple response to confirm scan started
type ScanResponse struct {
	Message    string `json:"message"`
	WorkflowID string `json:"workflow_id,omitempty"`
}

// AutoFixPlan is what Kestra will POST back to us after doing Syft + Grype + AI
// For now keep it simple – we can extend as we go.
type AutoFixPlan struct {
	RepoURL         string `json:"repo_url"`
	Summary         string `json:"summary"`
	FixStrategy     string `json:"fix_strategy"` // e.g. "bump_base_image"
	FromImage       string `json:"from_image,omitempty"`
	ToImage         string `json:"to_image,omitempty"`
	Vulnerabilities string `json:"vulnerability,omitempty"`
}

type StoredPlan struct {
	WorkflowID string      `json:"workflow_id"`
	Plan       AutoFixPlan `json:"plan"`
	Status     string      `json:"status"`
	FailedTask string      `json:"failed_task,omitempty"`
	Error      string      `json:"error"`
	CreatedAt  time.Time   `json:"created_at"`
}
