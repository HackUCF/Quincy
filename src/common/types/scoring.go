package types

// Score is the results of a completed score check.
// It has a boolean pass/fail, and a string containing an output message.
type Score struct {
	TeamNum     TeamNum     `json:"team_num"`
	Status      bool        `json:"status"`
	BoxName     BoxName     `json:"box"`
	ServiceName ServiceName `json:"service"`
	Message     string      `json:"message"`

	// Timestamp is inserted by the api on submission.
	// Does not need to added when building manually.
	Timestamp int64 `json:"timestamp"`
}

// ScoreResult is a flexible object used to store information about scoring results.
// This can be used for an entire team, one specific service for one team, one box for every team, and more.
type ScoreResult struct {
	ChecksPassed  uint64  `json:"checks_passed"`
	ChecksFailed  uint64  `json:"checks_failed"`
	TotalChecks   uint64  `json:"total_checks"`
	UptimePercent float64 `json:"uptime_percent"`
}
