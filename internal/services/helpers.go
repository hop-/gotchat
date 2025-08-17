package services

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type AtomicRunningStatus struct {
	runningStatus bool
	mu            sync.RWMutex
}

func (ars *AtomicRunningStatus) setRunningStatus(status bool) {
	ars.mu.Lock()
	defer ars.mu.Unlock()

	ars.runningStatus = status
}

func (ars *AtomicRunningStatus) isRunning() bool {
	ars.mu.RLock()
	defer ars.mu.RUnlock()

	return ars.runningStatus
}

func generateUuid() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return id.String()
}

func generateRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func generateRandomString() string {
	return fmt.Sprintf(
		"%d-%d-%d-%d",
		generateRandomInt(1000, 9999),
		generateRandomInt(1000, 9999),
		generateRandomInt(1000, 9999),
		generateRandomInt(1000, 9999),
	)
}
