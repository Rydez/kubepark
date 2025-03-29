package base

import (
	"kubepark/pkg/state"
)

// AttractionState represents the persistent state of an attraction
type AttractionState struct {
	IsPurchased bool `json:"is_purchased"`
}

// StateManager manages the attraction's persistent state
type StateManager struct {
	manager *state.Manager
}

// NewStateManager creates a new state manager
func NewStateManager(volumePath string) (*StateManager, error) {
	initialState := &AttractionState{
		IsPurchased: false,
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

// IsPurchased returns whether the attraction has been purchased
func (s *StateManager) IsPurchased() bool {
	return s.get().IsPurchased
}

// SetPurchased sets whether the attraction has been purchased
func (s *StateManager) SetPurchased(purchased bool) error {
	return s.set(func(state *AttractionState) {
		state.IsPurchased = purchased
	})
}
