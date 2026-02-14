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

// a runner contains all of the configuration for the agent.
type runner struct {
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
}

func (r *runner) run(id string) {
	// get next check from api server
	resp, err := http.Get(r.checkURL)
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
	scriptPath, err := r.getScript(check.CheckID)
	if err != nil {
		log.Error(
			"failed to get script path from check",
			"runner_id", id,
			"error", err,
			"check", check.CheckID,
		)
		return
	}

	out, err := r.runScript(scriptPath, check)
	if err != nil {
		log.Error(
			"failed to run check",
			"runner_id", id,
			"error", err,
			"check", check.CheckID,
		)
		return
	}

	// rc.messages <- fmt.Sprintf("err: %s", out.stderr.String())
	// rc.messages <- fmt.Sprintf("out: %s", out.stdout.String())
	// rc.messages <- fmt.Sprintf("%s", scriptPath)

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

	resp, err = http.Post(r.scoreURL, "application/json", bytes.NewReader(data))
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

func (r *runner) loop(id string) {
	for {
		time.Sleep(r.loopTime)
		r.run(id)
	}
}
