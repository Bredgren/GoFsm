/*
Package fsm provides basic types to implement state machines.
*/
package fsm

// State represents a state of the state machine. It provides a name which may be
// used to easily compare states (instead of by pointers or types).
type State interface {
	Name() string
	// OnEnter is called right after the state becomes the current state
	OnEnter()
	// OnExit is called right before the state changes
	OnExit()
}

type Fsm struct {
	CurrentState State
	PreviousState State
}

// NewFsm creates a new Fsm and in the given intial state. Causes the initial state's
// OnEnter function to be called. The Fsm's previousState is initially nil.
func NewFsm(initialState State) *Fsm {
	fsm := Fsm{initialState, nil}
	initialState.OnEnter()
	return &fsm
}

// Changes to the given state, triggering OnEnter and OnExit as appropriate.
func (f *Fsm) GotoState(state State) {
	f.CurrentState.OnExit()
	f.PreviousState = f.CurrentState
	f.CurrentState = state
	f.CurrentState.OnEnter()
}
