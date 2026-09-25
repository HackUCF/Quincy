package types

import "time"

// PauseState describes whether quincy should count scores during a period or not.
type PauseState bool

const (
	// Paused is the state where no scores are counted.
	Paused PauseState = true
	// Unpause is the state where scores are counted.
	Unpaused PauseState = false
	// PauseDefault can be used in error cases as a safe failure incase the caller ignores errors.
	PauseDefault PauseState = Paused
)

// Pause type defines the type of pause
type PauseType string

const (
	// ScoringPause is a type of pause where checks are still run, but scores are not counted
	// The scoreboard reflects the most recent checks, but final scores are not updated
	ScoringPause PauseType = "scoring"

	// CheckPause is a type of pause where the scoring agents do not run checks.
	CheckPause PauseType = "check"
)

// PauseRecord is the state of the competitions pauses.
type PauseRecord struct {
	ScoringState PauseState `json:"scoring_is_paused"   example:"false"`
	ScoringSince time.Time  `json:"scoring_since"       example:"2026-09-25T14:02:11.482913Z"`
	CheckState   PauseState `json:"checks_are_paused"   example:"true"`
	CheckSince   time.Time  `json:"checks_since"        example:"2026-09-25T14:37:50.119204Z"`
}
