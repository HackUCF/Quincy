package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HackUCF/quincy/common/types"

	_ "embed"
)

// a template passed as the message for a completed check.
//
//go:embed message_template.txt
var messageTpl string

// shared http client for requests to api
var apiClient http.Client

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

	// ensure it was recieved
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("bad status code posting request to quincy: %d", resp.StatusCode)
		return err
	}

	return nil
}
