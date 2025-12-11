package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /api/scan
// Body: { "repo_url": "https://github.com/user/repo.git" }
func handleScanRequest(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	// TODO: validate repo URL more strictly if needed

	// Trigger Kestra workflow here (Syft + Grype + AI)
	workflowID, err := triggerKestraWorkflow(req.RepoURL)
	if err != nil {
		log.Printf("failed to trigger Kestra workflow: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to start scan",
		})
		return
	}

	c.JSON(http.StatusAccepted, ScanResponse{
		Message:    "Scan started",
		WorkflowID: workflowID,
	})
}

// POST /api/autofix-plan
// Called by Kestra when it finishes AI planning
func handleAutoFixPlan(c *gin.Context) {
	var plan AutoFixPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid autofix plan payload",
		})
		return
	}

	log.Printf("Received AutoFix plan from Kestra: %+v\n", plan)

	// TODO:
	// 1. Clone repo
	// 2. Apply changes based on plan (update Dockerfile, etc.)
	// 3. Commit & push branch
	// 4. Open PR using GitHub API
	// For now, stub this:
	if err := applyAutoFixPlanAndOpenPR(plan); err != nil {
		log.Printf("failed to apply autofix plan: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to apply autofix plan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AutoFix plan applied and PR created",
	})
}
