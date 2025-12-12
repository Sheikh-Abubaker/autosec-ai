package main

import (
	"sync"
	"time"
)

var (
	planStore = make(map[string]StoredPlan)
	planMutex sync.RWMutex
)

func SavePlan(workflowID string, plan AutoFixPlan) {
	planMutex.Lock()
	defer planMutex.Unlock()
	planStore[workflowID] = StoredPlan{
		WorkflowID: workflowID,
		Plan:       plan,
		CreatedAt:  time.Now(),
	}
}

func GetPlan(workflowID string) (StoredPlan, bool) {
	planMutex.RLock()
	defer planMutex.RUnlock()
	p, ok := planStore[workflowID]
	return p, ok
}
