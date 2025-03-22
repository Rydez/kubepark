package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// GameState represents the current state of the park
type GameState struct {
	Cash        float64   `json:"cash"`
	CurrentTime time.Time `json:"current_time"`
	LastSaved   time.Time `json:"last_saved"`
	Mode        string    `json:"mode"`
	EntranceFee float64   `json:"entrance_fee"`
	OpensAt     int       `json:"opens_at"`
	ClosesAt    int       `json:"closes_at"`
	Closed      bool      `json:"closed"`
	VolumePath  string    `json:"-"`
	mu          sync.RWMutex
}

// New creates a new game state manager
func New(volumePath string) (*GameState, error) {
	state := &GameState{
		VolumePath:  volumePath,
		CurrentTime: time.Now(),
		LastSaved:   time.Now(),
		Cash:        1000, // Starting cash
	}

	// Try to load existing state if volume path is provided
	if volumePath != "" {
		if err := state.Load(); err != nil {
			// If file doesn't exist, that's fine - we'll create it on first save
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to load state: %w", err)
			}
		}
	}

	return state, nil
}

// Load loads the game state from the volume
func (s *GameState) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.VolumePath == "" {
		return nil
	}

	data, err := os.ReadFile(filepath.Join(s.VolumePath, "state.json"))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return nil
}

// Save saves the game state to the volume
func (s *GameState) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.VolumePath == "" {
		return nil
	}

	s.LastSaved = time.Now()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write to a temporary file first
	tmpFile := filepath.Join(s.VolumePath, "state.json.tmp")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Rename the temporary file to the actual file
	// This is an atomic operation that ensures we don't corrupt the state file
	if err := os.Rename(tmpFile, filepath.Join(s.VolumePath, "state.json")); err != nil {
		os.Remove(tmpFile) // Clean up the temporary file
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// AddCash adds to the park's cash amount
func (s *GameState) AddCash(amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Cash += amount
}

// SetCash sets the park's cash amount
func (s *GameState) SetCash(amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Cash = amount
}

// GetCash returns the park's current cash amount
func (s *GameState) GetCash() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Cash
}

// UpdateTime updates the park's current time
func (s *GameState) UpdateTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentTime = t
}

// GetTime returns the park's current time
func (s *GameState) GetTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.CurrentTime
}

// SetClosed sets whether the park is closed
func (s *GameState) SetClosed(closed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Closed = closed
}

// IsClosed returns whether the park is closed
func (s *GameState) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Closed
}
