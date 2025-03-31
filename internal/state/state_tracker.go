package state

import (
	"sync"
	"time"
)

type DeleteStatus string

const (
	StatusInitiated DeleteStatus = "initiated"
	StatusConfirmed DeleteStatus = "confirmed"
	StatusDeleting  DeleteStatus = "deleting"
	StatusDone      DeleteStatus = "done"
	StatusFailed    DeleteStatus = "failed"
)

type DeletionTracker struct {
	mu     sync.RWMutex
	status map[int]DeleteStatus
	expiry map[int]time.Time
}

// To keep it persistant
var tracker = &DeletionTracker{
	status: make(map[int]DeleteStatus),
	expiry: make(map[int]time.Time),
}

func SetStatus(evidenceID int, status DeleteStatus) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	tracker.status[evidenceID] = status
	tracker.expiry[evidenceID] = time.Now().Add(5 * time.Minute)
}

func GetStatus(evidenceID int) DeleteStatus {
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()
	if exp, exists := tracker.expiry[evidenceID]; !exists || time.Now().After(exp) {
		return ""
	}
	return tracker.status[evidenceID]
}

func ClearStatus(evidenceID int) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	delete(tracker.status, evidenceID)
	delete(tracker.expiry, evidenceID)
}
