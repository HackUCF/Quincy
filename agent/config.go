package agent

import (
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/HackUCF/Quincy/common/log"
)

// environment variables
const (
	// url of api server. can include protocol, port, and path
	// QU_API_URL='http://api.quin.cy'
	// QU_API_URL='https://proxy.quin.cy:8443/'
	envAPIURL = "QU_API_URL"

	// directory to store temporary files
	// these temp files are used to communicate with scripts
	// can be absolute or relative
	// defaults to 'tmp'
	envTempDir = "QU_TEMP_DIR"

	// directory to search for scripts
	// can be absolute or relative
	// defaults to 'scripts'
	envScriptsDir = "QU_CHECKS_DIR"

	// time in seconds to wait before fetching and running the next check.
	envLoopTime = "QU_LOOP_TIME"

	// number of concurrent scoring threads to run
	envNumThreads = "QU_NUM_THREADS"
)

// default values for environment variables
const (
	// default url for local api
	// should only be used for development
	defaultAPIURL = "http://127.0.0.1:8888"

	// directory used to generate temporary files for communication with scripts
	// platform independent ./tmp or .\tmp
	defaultTempDir = "tmp"

	// directory used to search for scripts
	// platform independent ./scripts or .\scripts
	defaultScriptsDir = "scripts"

	// a small default loop time for development
	defaultLoopTime = 1 * time.Second

	// a large default number of threads for development
	defaultNumThreads = 15
)

// contains all of the configuration for the agent.
type agentConfig struct {
	// url to get a new check from
	checkURL string

	// url to post a completed scorecheck to
	scoreURL string

	// directory to store temp files used to communcate with check scripts
	tempDir string

	// directory to find check scripts
	checksDir string

	// how long to sleep before looping
	loopTime time.Duration

	// how many agents to spawn
	numThreads int
}

// getConfig generates and returns the configuration object.
// It pulls from environment variables and does some construction.
// It will panic if anything here fails.
func getConfig() agentConfig {
	var cfg agentConfig

	// load the loop time env var
	// ugly
	strLoopTime := os.Getenv(envLoopTime)
	if strLoopTime == "" {
		cfg.loopTime = defaultLoopTime
	} else {
		// convert to time duration if it exists
		intLoopTime, err := strconv.Atoi(strLoopTime)

		if err != nil {
			// check that the env var doesn't have nonsense in it
			log.Warn(
				"loop time environment variable could not be parsed as an integer, using the default.",
				"var", envLoopTime,
				"value", strLoopTime,
				"default", defaultLoopTime,
			)
			cfg.loopTime = defaultLoopTime
		} else {
			cfg.loopTime = time.Duration(intLoopTime) * time.Second
		}
	}

	// load the number of threads
	strNumThreads := os.Getenv(envNumThreads)
	if strNumThreads == "" {
		cfg.numThreads = defaultNumThreads
	} else {
		// convert to integer
		intNumThreads, err := strconv.Atoi(strNumThreads)

		if err != nil {
			// check that the env var doesn't have nonsense in it
			log.Warn(
				"num threads environment variable could not be parsed as an integer, using the default.",
				"var", envNumThreads,
				"value", strNumThreads,
				"default", defaultNumThreads,
			)
			cfg.numThreads = defaultNumThreads
		} else {
			cfg.numThreads = intNumThreads
		}
	}

	relTempDir := os.Getenv(envTempDir)
	if relTempDir == "" {
		relTempDir = defaultTempDir
	}

	var err error
	cfg.tempDir, err = filepath.Abs(relTempDir)
	if err != nil {
		log.Panic(
			"failed to create absolute path from relative temporary directory",
			"temp_dir", relTempDir,
			"error", err,
		)
	}

	relChecksDir := os.Getenv(envScriptsDir)
	if relChecksDir == "" {
		relChecksDir = defaultScriptsDir
	}

	cfg.checksDir, err = filepath.Abs(relChecksDir)
	if err != nil {
		log.Panic(
			"failed to create absolute path from relative checks directory",
			"checks_dir", relChecksDir,
			"error", err,
		)
	}

	apiURL := os.Getenv(envAPIURL)
	if apiURL == "" {
		apiURL = defaultAPIURL

		log.Warn(
			"api url is not set, using the default",
			"api_url", apiURL,
		)
	}

	checksURL := "/api/v1/checks"
	cfg.checkURL, err = url.JoinPath(apiURL, checksURL)
	if err != nil {
		log.Panic(
			"failed to construct checks url endpoint",
			"base", apiURL,
			"path", checksURL,
		)
	}

	scoresURL := "/api/v1/scores"
	cfg.scoreURL, err = url.JoinPath(apiURL, scoresURL)
	if err != nil {
		log.Panic(
			"failed to construct scores url endpoint",
			"base", apiURL,
			"path", scoresURL,
		)
	}

	return cfg
}
