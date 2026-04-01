package agent

import (
	"net/url"
	"path/filepath"
	"time"

	"github.com/HackUCF/Quincy/common/log"
)

// Config holds the viper-loadable agent configuration.
// Fields map to viper keys, which in turn map to QU_<KEY> env vars
// and --<key> CLI flags (e.g. api_url → QU_API_URL / --api-url).
type Config struct {
	APIURL     string `mapstructure:"api_url"`
	ChecksDir  string `mapstructure:"checks_dir"`
	LoopTime   int    `mapstructure:"loop_time"` // seconds
	NumThreads int    `mapstructure:"num_threads"`
}

// agentConfig is the resolved runtime config derived from Config.
// It holds constructed URLs, an absolute path, and a time.Duration.
type agentConfig struct {
	checkURL   string
	scoreURL   string
	checksDir  string
	loopTime   time.Duration
	numThreads int
}

// resolve derives the runtime agentConfig from the raw Config.
func (cfg *Config) resolve() *agentConfig {
	var ac agentConfig
	ac.loopTime = time.Duration(cfg.LoopTime) * time.Second
	ac.numThreads = cfg.NumThreads

	var err error
	ac.checksDir, err = filepath.Abs(cfg.ChecksDir)
	if err != nil {
		log.Panic("failed to resolve checks directory", "checks_dir", cfg.ChecksDir, "error", err)
	}

	ac.checkURL, err = url.JoinPath(cfg.APIURL, "/api/v1/agent/new-check")
	if err != nil {
		log.Panic("failed to construct check URL", "api_url", cfg.APIURL, "error", err)
	}

	ac.scoreURL, err = url.JoinPath(cfg.APIURL, "/api/v1/agent/completed-score")
	if err != nil {
		log.Panic("failed to construct score URL", "api_url", cfg.APIURL, "error", err)
	}

	return &ac
}
