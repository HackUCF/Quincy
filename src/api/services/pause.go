package services

import (
	"sync/atomic"
)

var isPaused atomic.Bool

func Pause() {
	isPaused.Store(true)
}

func Unpause() {
	isPaused.Store(false)
}

func IsPaused() bool { return isPaused.Load() }
