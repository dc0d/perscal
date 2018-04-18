package main

//-----------------------------------------------------------------------------

// activate activates the (start) state and it's consecutive
// states until the next state is nil or encounters an error
func activate(s state) (_err error) {
	next := s
	for next != nil && _err == nil {
		next, _err = next()
	}
	return
}

// state is a function that represents a State
type state func() (state, error)

//-----------------------------------------------------------------------------
