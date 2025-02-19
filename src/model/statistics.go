package model

import (
	"sync"
)

// Statistics tracks various account processing statistics
type Statistics struct {
	mu sync.Mutex // Protects all fields

	total     int // Total number of accounts
	processed int // Number of accounts processed
	success   int // Number of successful accounts
	failed    int // Number of failed accounts
	locked    int // Number of locked accounts
	suspended int // Number of suspended accounts
}

// NewStatistics creates a new Statistics instance
func NewStatistics(total int) *Statistics {
	return &Statistics{
		total: total,
	}
}

// GetStats returns all statistics values
func (s *Statistics) GetStats() (total, processed, success, failed, locked, suspended int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total, s.processed, s.success, s.failed, s.locked, s.suspended
}

// IncrementProcessed increments the processed count
func (s *Statistics) IncrementProcessed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed++
}

// IncrementSuccess increments the success count
func (s *Statistics) IncrementSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.success++
}

// IncrementFailed increments the failed count
func (s *Statistics) IncrementFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failed++
}

// IncrementLocked increments the locked count
func (s *Statistics) IncrementLocked() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.locked++
}

// IncrementSuspended increments the suspended count
func (s *Statistics) IncrementSuspended() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.suspended++
}
