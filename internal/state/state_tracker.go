package state

import (
	"strconv"
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
	status map[string]DeleteStatus
	expiry map[string]time.Time
}

var tracker = &DeletionTracker{
	status: make(map[string]DeleteStatus),
	expiry: make(map[string]time.Time),
}

func makeKey(userID string, evidenceID int) string {
	return userID + ":" + strconv.Itoa(evidenceID)
}

func SetStatus(userID string, evidenceID int, status DeleteStatus) {
	key := makeKey(userID, evidenceID)
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	tracker.status[key] = status
	tracker.expiry[key] = time.Now().Add(5 * time.Minute)
}

func GetStatus(userID string, evidenceID int) DeleteStatus {
	key := makeKey(userID, evidenceID)
	tracker.mu.RLock()
	defer tracker.mu.RUnlock()
	if time.Now().After(tracker.expiry[key]) {
		return ""
	}
	return tracker.status[key]
}

func ClearStatus(userID string, evidenceID int) {
	key := makeKey(userID, evidenceID)
	tracker.mu.Lock()
	defer tracker.mu.Unlock()
	delete(tracker.status, key)
	delete(tracker.expiry, key)
}
