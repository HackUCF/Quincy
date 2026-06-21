package types

// Score is the results of a completed score check.
// It has a boolean pass/fail, and a string containing an output message.
type Score struct {
	TeamNum     TeamNum     `json:"team_num"  example:"1"`
	Status      bool        `json:"status"    example:"true"`
	BoxName     BoxName     `json:"box"       example:"scrapyard"`
	ServiceName ServiceName `json:"service"   example:"blog"`
	Message     string      `json:"message"   example:"HTTP 200 OK\nexit 0"`

	// Timestamp is inserted by the api on submission.
	// Does not need to added when building manually.
	Timestamp int64 `json:"timestamp" example:"1700000000000000"`
}

// ScoreResult is a flexible object used to store information about scoring results.
// This can be used for an entire team, one specific service for one team, one box for every team, and more.
type ScoreResult struct {
	ChecksPassed  uint64  `json:"checks_passed"  example:"42"`
	ChecksFailed  uint64  `json:"checks_failed"  example:"8"`
	TotalChecks   uint64  `json:"total_checks"   example:"50"`
	UptimePercent float64 `json:"uptime_percent" example:"84.00"`
}
