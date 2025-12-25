package main

import "github.com/gin-gonic/gin"

func registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/scan", handleScanRequest)
		api.POST("/autofix-plan", handleAutoFixPlan) // callback from Kestra
		api.POST("/scan-failure", handleScanFailure) // NEW
		// NEW ENDPOINT
		api.GET("/scan/:workflow_id", handleGetScanStatus)
	}
}
