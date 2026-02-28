package graphs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"slices"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"
)

// scoreboardPoint is the a specific check on the scoreboard chart.
// this is the format the js library expects
type scoreboardPoint struct {
	X         string `json:"x"`
	Y         string `json:"y"`
	V         int    `json:"v"`
	Timestamp string `json:"ts"`
	Message   string `json:"msg"`
}

// rawScoreboardData contains the data required to render a scoreboard chart.
// it's not rendered yet.
type rawScoreboardData struct {
	XLabels    []string
	YLabels    []string
	DataPoints []scoreboardPoint
}

// ScoreboardData contains the rendered string required for the scoreboard graph template.
type ScoreboardData struct {
	// list of "boxid-serviceid"
	XLabels template.JS
	// list of team labels "team xx"
	YLabels template.JS
	// the datapoints, a list of objects containing x, y, and some values.
	DataPoints template.JS
}

// GetScoreboardData returns the template data needed to render a scoreboard chart.
// data contains rendered json strings.
func GetScoreboardData() (*ScoreboardData, error) {
	db := conn.Get()

	rows, err := db.Query(`
		SELECT service, box, team_num, status, timestamp, message
    FROM recent_scores
    ORDER BY team_num, box, service
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := new(rawScoreboardData)
	data.XLabels = make([]string, 0)
	data.YLabels = make([]string, 0)
	data.DataPoints = make([]scoreboardPoint, 0)

	for rows.Next() {
		// scan the necessary scores out
		var s types.Score
		rows.Scan(&s.ServiceID, &s.BoxID, &s.TeamNum, &s.Status, &s.Timestamp, &s.Message)

		// get x, y, and v
		var point scoreboardPoint
		point.X = string(s.BoxID) + "-" + string(s.ServiceID)
		point.Y = "team " + fmt.Sprintf("%d", s.TeamNum)
		if s.Status {
			point.V = 1
		} else {
			point.V = 0
		}
		point.Timestamp = fmt.Sprintf("%d", s.Timestamp)
		point.Message = s.Message

		// add to labels if not there already
		if !slices.Contains(data.XLabels, point.X) {
			data.XLabels = append(data.XLabels, point.X)
		}
		if !slices.Contains(data.YLabels, point.Y) {
			data.YLabels = append(data.YLabels, point.Y)
		}

		// add to list of full data
		data.DataPoints = append(data.DataPoints, point)
	}

	// dump to json
	xBytes, err := json.Marshal(data.XLabels)
	if err != nil {
		return nil, err
	}
	yBytes, err := json.Marshal(data.YLabels)
	if err != nil {
		return nil, err
	}
	datasetBytes, err := json.Marshal(data.DataPoints)
	if err != nil {
		return nil, err
	}

	// return
	retval := new(ScoreboardData)
	retval.XLabels = template.JS(xBytes)
	retval.YLabels = template.JS(yBytes)
	retval.DataPoints = template.JS(datasetBytes)

	log.Info(string(retval.XLabels))
	log.Info(string(retval.YLabels))
	log.Info(string(retval.DataPoints))

	return retval, nil
}
