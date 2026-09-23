package services

import (
	"sync/atomic"
	"time"

	"github.com/HackUCF/Quincy/src/common/log"
)

type pauseStateType struct {
	paused bool
	since  time.Time
}

var (
	// manages pause status safely from multiple threads
	// more complicated than a lock, but much faster
	pauseState atomic.Pointer[pauseStateType]
)

// IsPaused quickly checks if quincy is paused.
// True means paused, false means unpaused.
func IsPaused() bool { return pauseState.Load().paused }

func setPauseState(desiredState bool) (changed bool) {

	// loop until no contention
	// this basically never loops unless someone's spamming the pause or unpause endpoints
	// this just guarantees safety
	for {

		// load the old state
		old := pauseState.Load()
		if old.paused == desiredState {
			// early return if already the desired state
			return false
		}

		// make the new state
		new := &pauseStateType{paused: desiredState, since: time.Now()}

		// swap only if no one contented since loading
		swapSucceeded := pauseState.CompareAndSwap(old, new)

		if swapSucceeded {

			// log and exit when state changes
			log.Info(
				"quincy pause state changed",
				"old_state", old.paused,
				"new_state", new.paused,
				"old_since", old.since,
				"new_since", new.since,
			)
			return true
		}
	}
}

// Pause pauses quincy, updating the time it entered this state only if the state changes.
func Pause() (changed bool) { return setPauseState(true) }

// Unpause unpauses quincy, updating the stored time changed only if the state changes.
func Unpause() (changed bool) { return setPauseState(false) }

// PauseStatus returns the competition pause status as well as the time it entered that state.
func PauseStatus() (isPaused bool, since time.Time) {
	state := pauseState.Load()
	return state.paused, state.since
}
