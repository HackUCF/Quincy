package main

import (
	"net/url"
	"os"
	"path/filepath"
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
)

// default values for environment variables
const (
	// default url for local api
	// should only be uised for development
	defaultAPIURL = "http://127.0.0.1:8888"

	// directory used to generate temporary files for communication with scripts
	// platform independent ./tmp or .\tmp
	defaultTempDir = "tmp"

	// directory used to search for scripts
	// platform independent ./scripts or .\scripts
	defaultScriptsDir = "scripts"
)

// getRunner calculates and returns the runner configuration object.
// It pulls from environment variables and does some construction.
// It will panic if anything here fails.
func getRunner() runner {
	var r runner

	r.loopTime = 3 * time.Second

	relTempDir := os.Getenv(envTempDir)
	if relTempDir == "" {
		relTempDir = defaultTempDir
	}

	var err error
	r.tempDir, err = filepath.Abs(relTempDir)
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

	r.checksDir, err = filepath.Abs(relChecksDir)
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
	r.checkURL, err = url.JoinPath(apiURL, checksURL)
	if err != nil {
		log.Panic(
			"failed to construct checks url endpoint",
			"base", apiURL,
			"path", checksURL,
		)
	}

	scoresURL := "/api/v1/scores"
	r.scoreURL, err = url.JoinPath(apiURL, scoresURL)
	if err != nil {
		log.Panic(
			"failed to construct scores url endpoint",
			"base", apiURL,
			"path", scoresURL,
		)
	}

	return r
}
