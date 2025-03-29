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
	LastSaved   time.Time `json:"last_saved"`
	Mode        string    `json:"mode"`
	EntranceFee float64   `json:"entrance_fee"`
	OpensAt     int       `json:"opens_at"`
	ClosesAt    int       `json:"closes_at"`
	Closed      bool      `json:"closed"`
	TotalSpace  float64   `json:"total_space"` // Total park space in acres
	UsedSpace   float64   `json:"used_space"`  // Used space in acres
	GuestSize   float64   `json:"guest_size"`  // Size of each guest in acres
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
		OpensAt:     config.OpensAt,
		ClosesAt:    config.ClosesAt,
		Closed:      config.Closed,
		GuestSize:   0.1, // Each guest takes 0.1 acres
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

// SetClosed sets whether the park is closed
func (s *StateManager) SetClosed(closed bool) error {
	return s.set(func(state *ParkState) {
		state.Closed = closed
	})
}

// IsClosed returns whether the park is closed
func (s *StateManager) IsClosed() bool {
	return s.get().Closed
}

// CanAddGuest checks if there's enough space for a new guest
func (s *StateManager) CanAddGuest() bool {
	usedSpace := s.get().UsedSpace
	guestSize := s.get().GuestSize
	totalSpace := s.get().TotalSpace
	return usedSpace+guestSize <= totalSpace
}

// AddGuest adds a guest to the park
func (s *StateManager) AddGuest() error {
	if !s.CanAddGuest() {
		return fmt.Errorf("not enough space in park for guest. Need %.1f acres but only have %.1f acres available", s.get().GuestSize, s.get().TotalSpace-s.get().UsedSpace)
	}

	return s.set(func(state *ParkState) {
		state.UsedSpace += state.GuestSize
	})
}

// RemoveGuest removes a guest from the park
func (s *StateManager) RemoveGuest() error {
	return s.set(func(state *ParkState) {
		state.UsedSpace -= state.GuestSize
	})
}

// GetAvailableSpace returns the amount of space available in the park
func (s *StateManager) GetAvailableSpace() float64 {
	return s.get().TotalSpace - s.get().UsedSpace
}

// GetTotalSpace returns the total space in the park
func (s *StateManager) GetTotalSpace() float64 {
	return s.get().TotalSpace
}
