package service

import (
	"log"
	"time"

	"github.com/rnegic/synchronous/internal/entity"
	"github.com/rnegic/synchronous/internal/interfaces"
)

// SessionCleanupService handles automatic cleanup of stale sessions
type SessionCleanupService struct {
	sessionRepo interfaces.SessionRepository
	interval    time.Duration
	maxAge      time.Duration
}

// NewSessionCleanupService creates a new cleanup service
func NewSessionCleanupService(
	sessionRepo interfaces.SessionRepository,
	interval time.Duration,
	maxAge time.Duration,
) *SessionCleanupService {
	return &SessionCleanupService{
		sessionRepo: sessionRepo,
		interval:    interval,
		maxAge:      maxAge,
	}
}

// Start begins the cleanup routine
func (s *SessionCleanupService) Start() {
	log.Printf("[SessionCleanup] 🧹 Starting cleanup service (interval: %v, maxAge: %v)\n", s.interval, s.maxAge)

	ticker := time.NewTicker(s.interval)
	go func() {
		for range ticker.C {
			s.cleanup()
		}
	}()
}

// cleanup removes stale pending sessions
func (s *SessionCleanupService) cleanup() {
	sessions, err := s.sessionRepo.GetSessionsByStatus(entity.SessionStatusPending)
	if err != nil {
		log.Printf("[SessionCleanup] ❌ Failed to get pending sessions: %v\n", err)
		return
	}

	now := time.Now()
	cleaned := 0

	for _, session := range sessions {
		age := now.Sub(session.CreatedAt)
		if age > s.maxAge {
			// Mark session as completed (or you could delete it)
			completedAt := now
			session.Status = entity.SessionStatusCompleted
			session.CompletedAt = &completedAt

			if err := s.sessionRepo.Update(session); err != nil {
				log.Printf("[SessionCleanup] ❌ Failed to cleanup session %s: %v\n", session.ID, err)
				continue
			}

			cleaned++
			log.Printf("[SessionCleanup] 🧹 Cleaned up stale session: %s (age: %v)\n", session.ID, age)
		}
	}

	if cleaned > 0 {
		log.Printf("[SessionCleanup] ✅ Cleaned up %d stale sessions\n", cleaned)
	}
}
