package main

import (
	"fmt"
	"kubepark/pkg/state"
	"time"
)

// ParkState represents the persistent state of the park
type ParkState struct {
	Money       float64   `json:"money"`
	CurrentTime time.Time `json:"current_time"`
	Mode        string    `json:"mode"`
	EntranceFee float64   `json:"entrance_fee"`
	TotalSpace  float64   `json:"total_space"` // Total park space in acres
}

// StateManager manages the attraction's persistent state
type StateManager struct {
	manager *state.Manager
}

// NewStateManager creates a new state manager
func NewStateManager(config *Config) (*StateManager, error) {
	initialState := &ParkState{
		Money:       100000, // Start with $100,000
		CurrentTime: time.Now(),
		Mode:        config.Mode,
		EntranceFee: config.EntranceFee,
	}

	// Set total space based on mode
	switch initialState.Mode {
	case "easy":
		initialState.TotalSpace = 300
	case "medium":
		initialState.TotalSpace = 100
	case "hard":
		initialState.TotalSpace = 10
	default:
		return nil, fmt.Errorf("mode not set on park")
	}

	manager, err := state.New(initialState, config.VolumePath)
	if err != nil {
		return nil, err
	}

	return &StateManager{
		manager: manager,
	}, nil
}

func (s *StateManager) set(setter func(*ParkState)) error {
	state := s.manager.Get().(*ParkState)
	setter(state)
	return s.manager.Set(state)
}

func (s *StateManager) get() *ParkState {
	return s.manager.Get().(*ParkState)
}

func (s *StateManager) GetEntranceFee() float64 {
	return s.get().EntranceFee
}

// AddCash adds to the park's cash amount
func (s *StateManager) AddMoney(amount float64) error {
	return s.set(func(state *ParkState) {
		state.Money += amount
	})
}

// SetCash sets the park's cash amount
func (s *StateManager) SetMoney(amount float64) error {
	return s.set(func(state *ParkState) {
		state.Money = amount
	})
}

// GetCash returns the park's current cash amount
func (s *StateManager) GetMoney() float64 {
	return s.get().Money
}

// GetTime returns the park's current time
func (s *StateManager) GetTime() time.Time {
	return s.get().CurrentTime
}

// SetTime sets the park's current time
func (s *StateManager) SetTime(t time.Time) error {
	return s.set(func(state *ParkState) {
		state.CurrentTime = t
	})
}

// GetTotalSpace returns the total space in the park
func (s *StateManager) GetTotalSpace() float64 {
	return s.get().TotalSpace
}
