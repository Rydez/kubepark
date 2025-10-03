package base

import (
	"kubepark/pkg/state"
)

// AttractionState represents the persistent state of an attraction
type AttractionState struct {
	IsBroken bool `json:"is_broken"`
}

// StateManager manages the attraction's persistent state
type StateManager struct {
	manager *state.Manager
}

// NewStateManager creates a new state manager
func NewStateManager(volumePath string) (*StateManager, error) {
	initialState := &AttractionState{
		IsBroken: false,
	}

	manager, err := state.New(initialState, volumePath)
	if err != nil {
		return nil, err
	}

	return &StateManager{
		manager: manager,
	}, nil
}

func (s *StateManager) set(setter func(*AttractionState)) error {
	state := s.manager.Get().(*AttractionState)
	setter(state)
	return s.manager.Set(state)
}

func (s *StateManager) get() *AttractionState {
	return s.manager.Get().(*AttractionState)
}

// IsBroken returns whether the attraction is broken
func (s *StateManager) IsBroken() bool {
	return s.get().IsBroken
}

// SetBroken sets whether the attraction is broken
func (s *StateManager) SetBroken(broken bool) error {
	return s.set(func(state *AttractionState) {
		state.IsBroken = broken
	})
}
