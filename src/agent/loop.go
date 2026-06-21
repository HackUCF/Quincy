package agent

import (
	"net/url"
	"time"

	"github.com/HackUCF/quincy/common/log"
)

func (cfg *AgentConfig) Loop() {

	// precompute some urls to save time later:
	serviceURL, err := url.JoinPath(cfg.APIURL, "/api/v1/agent/new-check")
	if err != nil {
		log.Panic("failed to construct check URL", "api_url", cfg.APIURL, "error", err)
	}
	scoreURL, err := url.JoinPath(cfg.APIURL, "/api/v1/agent/completed-score")
	if err != nil {
		log.Panic("failed to construct score URL", "api_url", cfg.APIURL, "error", err)
	}

	// convert configured loop time into seconds
	loopTime := time.Second * time.Duration(cfg.LoopTime)

	// infinite loop
	for {

		// start with sleep so continues don't skip it
		time.Sleep(loopTime)

		start := time.Now()

		// get service from quincy
		svc, err := getService(serviceURL)
		if err != nil {
			log.Error(
				"failed to get service from server",
				"error", err,
				"url", serviceURL,
			)
			continue
		}

		// run it (with the configured timeout)
		timeout := getTimeout(cfg, svc)
		output, err := runCheck(svc, timeout)
		if err != nil {
			log.Error(
				"failed to run check",
				"error", err,
				"service", svc,
			)
			continue
		}

		// convert output to api object
		score := makeScore(svc, output)

		// post result to quincy
		err = postScore(scoreURL, score)
		if err != nil {
			log.Error(
				"failed to post completed score to server",
				"error", err,
				"url", scoreURL,
				"score", score,
			)
			continue
		}

		// log success then loop again!
		dur := time.Since(start)
		log.Info(
			"successfully posted score to server",
			"duration", dur.Seconds(),
		)
	}
}
