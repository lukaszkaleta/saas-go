package universal

import "fmt"

type State interface {
	Name() string
	Change(newState string) error
}
type SolidState struct {
	State     State
	Current   string
	Available []string
}

func NewSolidState(current string, available []string, state State) *SolidState {
	return &SolidState{
		State:     state,
		Current:   current,
		Available: available,
	}
}

func (s *SolidState) Name() string {
	return s.Current
}

func (s *SolidState) Change(newState string) error {
	// ensure the requested state is one of the available states
	for _, a := range s.Available {
		if a == newState {
			s.Current = newState
			return s.State.Change(newState)
		}
	}
	return fmt.Errorf("state '%s' is not available", newState)
}
