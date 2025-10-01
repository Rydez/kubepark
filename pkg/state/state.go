package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Manager manages persistent state
type Manager struct {
	state      interface{}
	volumePath string
	mu         sync.RWMutex
}

// New creates a new state manager
func New(initialState interface{}, volumePath string) (*Manager, error) {
	manager := &Manager{
		state:      initialState,
		volumePath: volumePath,
	}

	// Load existing state if available
	if err := manager.Load(); err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	return manager, nil
}

// Load loads the state from disk
func (s *Manager) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.volumePath == "" {
		return nil
	}

	data, err := os.ReadFile(filepath.Join(s.volumePath, "state.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing state file, use defaults
		}
		return fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, s.state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return nil
}

// Save saves the state to disk
func (s *Manager) Save() error {
	if s.volumePath == "" {
		return nil
	}

	data, err := json.MarshalIndent(s.state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Write to a temporary file first
	tmpFile := filepath.Join(s.volumePath, "state.json.tmp")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Rename the temporary file to the actual file
	// This is an atomic operation that ensures we don't corrupt the state file
	if err := os.Rename(tmpFile, filepath.Join(s.volumePath, "state.json")); err != nil {
		os.Remove(tmpFile) // Clean up the temporary file
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// Get returns the current state
func (s *Manager) Get() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// Set updates the state and saves it to disk
func (s *Manager) Set(newState interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = newState
	return s.Save()
}
