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

func handleGetScanStatus(c *gin.Context) {
	workflowID := c.Param("workflow_id")

	planData, ok := GetPlan(workflowID)
	if !ok {
		c.JSON(200, gin.H{
			"status": "running",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "done",
		"plan":   planData.Plan,
	})
}

// POST /api/autofix-plan
// Called by Kestra when it finishes AI planning
func handleAutoFixPlan(c *gin.Context) {
	var plan AutoFixPlan

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(400, gin.H{"error": "invalid autofix payload"})
		return
	}

	log.Printf("Received AutoFix plan from Kestra: %+v\n", plan)

	workflowID := c.GetHeader("X-Kestra-Execution-Id")
	if workflowID == "" {
		workflowID = c.Query("workflow_id")
	}

	if workflowID == "" {
		c.JSON(400, gin.H{"error": "missing workflow_id"})
		return
	}

	SavePlan(workflowID, plan)

	log.Printf("AutoFix plan stored for workflow %s: %+v", workflowID, plan)

	c.JSON(200, gin.H{
		"message":     "plan stored",
		"workflow_id": workflowID,
	})
}
