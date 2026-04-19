package agent

import "strings"

// AgentConfig holds the viper-loadable agent configuration.
// Fields map to viper keys, which in turn map to QU_<KEY> env vars
// and --<key> CLI flags (e.g. api_url → QU_API_URL / --api-url).
type AgentConfig struct {
	APIURL     string `mapstructure:"api_url"`
	LoopTime   int    `mapstructure:"loop_time"` // seconds
	NumThreads int    `mapstructure:"num_threads"`
}

const Pass bool = true
const Fail bool = false

type scriptOutput struct {
	Status bool
	Stdout strings.Builder
	Stderr strings.Builder
}
