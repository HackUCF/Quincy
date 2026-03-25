package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"
)

// global http client object
var client http.Client

// run a single check.
// uses info from the config to determin behavior.
func (cfg *agentConfig) run(id string) {

	// get next check from api server
	resp, err := client.Get(cfg.checkURL)
	if err != nil {
		log.Error(
			"failed to get a new request from api server",
			"runner_id", id,
			"error", err,
		)
		return
	}

	// load service information from response body
	var check types.Service
	err = json.NewDecoder(resp.Body).Decode(&check)
	if err != nil {
		log.Error(
			"failed to load check json from response body",
			"runner_id", id,
			"error", err,
		)
		return
	}

	err = resp.Body.Close()
	if err != nil {
		log.Error(
			"failed to close check response body. this should never happen",
			"runner_id", id,
			"error", err,
		)
		return
	}

	// identify script from check id
	scriptPath, err := cfg.getScript(check.CheckID)
	if err != nil {
		log.Error(
			"failed to get script path from check",
			"runner_id", id,
			"error", err,
			"check", check.CheckID,
		)
		return
	}

	out, err := cfg.runScript(scriptPath, check)
	if err != nil {
		log.Error(
			"failed to run check",
			"runner_id", id,
			"error", err,
			"check", check.CheckID,
		)
		return
	}

	score := types.Score{
		TeamNum:   check.TeamNum,
		Status:    out.status,
		BoxID:     check.BoxID,
		ServiceID: check.ID,
		Message:   fmt.Sprintf("err %s\nout %s", out.stderr.String(), out.stdout.String()),
	}

	data, err := json.Marshal(score)
	if err != nil {
		log.Error(
			"failed to marshall score result into json",
			"runner_id", id,
			"error", err,
		)
		return
	}

	resp, err = client.Post(cfg.scoreURL, "application/json", bytes.NewReader(data))
	if err != nil {
		log.Error(
			"failed to post completed scorecheck to api server",
			"runner_id", id,
			"error", err,
		)
		return
	}

	err = resp.Body.Close()
	if err != nil {
		log.Error(
			"failed to close response body. this should never happen.",
			"runner_id", id,
			"error", err,
		)
		return
	}
}

// infinite loop of checks
func (cfg *agentConfig) loop(id string) {
	for {
		cfg.run(id)
		time.Sleep(cfg.loopTime)
	}
}
