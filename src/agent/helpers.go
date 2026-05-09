package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/HackUCF/quincy/common/types"

	_ "embed"
)

// a template passed as the message for a completed check.
//
//go:embed message_template.txt
var messageTpl string

// shared http client for requests to api
var apiTransport = &http.Transport{
	MaxIdleConns:        1000,
	MaxIdleConnsPerHost: 1000,
	MaxConnsPerHost:     0,
	IdleConnTimeout:     90 * time.Second,
}
var apiClient = &http.Client{Transport: apiTransport}

// grab next service to check from api
func getService(url string) (*types.Service, error) {

	// return value
	var svc types.Service

	// make http request
	resp, err := apiClient.Get(url)
	if err != nil {
		err := fmt.Errorf("failed to request check from quincy: %w", err)
		return nil, err
	}

	// load json from response
	err = json.NewDecoder(resp.Body).Decode(&svc)
	defer resp.Body.Close()
	if err != nil {
		err := fmt.Errorf("failed load quincy response as json: %w", err)
		return nil, err
	}

	return &svc, nil
}

// get the timeout from config or service.
// if the service-specific timeout is specified, use it.
// if the loop time is reasonable (>5s), use it.
// otherwise just pick a chill 30s.
func getTimeout(cfg *AgentConfig, svc *types.Service) time.Duration {

	var timeout time.Duration

	if svc.Timeout != 0 {
		// use the service-specific timeout if provided
		timeout = time.Duration(svc.Timeout * float64(time.Second))
	} else if cfg.LoopTime >= 5 {
		// use the loop time if long enough
		timeout = time.Duration(cfg.LoopTime) * time.Second
	} else {
		// default to 30s (unlikely to trigger)
		timeout = 30 * time.Second
	}

	return timeout
}

// convert the output into api compatible object
func makeScore(svc *types.Service, output *scriptOutput) *types.Score {
	score := new(types.Score)

	score.BoxName = svc.BoxName
	score.ServiceName = svc.Name
	score.TeamNum = svc.TeamNum
	score.Status = output.Status
	score.Message = fmt.Sprintf(
		messageTpl,
		output.Stdout.String(),
		output.Stderr.String(),
	)

	return score
}

// post result back to quincy
func postScore(url string, score *types.Score) error {

	// dump to json bytes
	obj, err := json.Marshal(score)
	if err != nil {
		err := fmt.Errorf("failed to dump score to json: %w", err)
		return err
	}

	// post to server
	resp, err := apiClient.Post(url, "application/json", bytes.NewReader(obj))
	if err != nil {
		err := fmt.Errorf("failed to request check from quincy: %w", err)
		return err
	}
	defer resp.Body.Close()

	// ensure it was recieved
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("bad status code posting request to quincy: %d", resp.StatusCode)
		return err
	}

	return nil
}
