package graphs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"slices"
	"time"

	"github.com/HackUCF/Quincy/api/db/misc"
	"github.com/HackUCF/Quincy/common/types"
)

type scoresDataset struct {
	Label string   `json:"label"`
	Data  []uint64 `json:"data"`
}

type rawScoresData struct {
	Labels []string
	Data   []scoresDataset
}

// ScoresData returns the template data needed to render the final scores length
type ScoresData struct {
	Labels   template.JS
	Datasets template.JS
}

// name format of each line, %d is the team number
const labelFmt = "team %d"

// how many samples to make for the graph.
// this is not an exact number.
// pauses and other gaps will mean this is a ceiling that almost never gets hit.
const numIntervals = 100

// GetScoresData contains the data required to render a line chart with all team scores.
func GetScoresData(db *sql.DB) (*ScoresData, error) {

	// get microsecond intervals for this competition
	// compeition duration divided, in microseconds
	duration, err := misc.GetCompDuration(db)
	if err != nil {
		return nil, err
	}
	bucketUs := (duration / numIntervals).Truncate(time.Second).Microseconds()

	// contains evil integer rounding hack. i do not like this
	rows, err := db.Query(`
			SELECT team_num, (timestamp / ?) * ? as bucket, SUM(status)
			FROM scores
			GROUP BY team_num, bucket
			ORDER BY team_num, bucket;
		`, bucketUs, bucketUs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// full data container setup
	var data rawScoresData
	data.Labels = make([]string, 0)
	data.Data = make([]scoresDataset, 0)

	// loop variables
	currentTeam := types.TeamNum(1) // messy
	var currentDataset scoresDataset
	currentDataset.Label = fmt.Sprintf(labelFmt, 1)
	currentDataset.Data = make([]uint64, 0)
	currentPoints := uint64(0)

	for rows.Next() {
		// temp variables
		var team types.TeamNum
		var timestamp int64
		var points uint64

		// scan into variables
		err := rows.Scan(&team, &timestamp, &points)
		if err != nil {
			return nil, err
		}

		// convert ts to string
		timeString := time.Unix(timestamp/1_000_000, 0).Format(time.TimeOnly)

		// add ts to graph labels (if missing)
		if !slices.Contains(data.Labels, timeString) {
			data.Labels = append(data.Labels, timeString)
		}

		// hacky
		// "we are done with this team, move on to the next"
		if team != currentTeam {
			// add completed dataset to final collection
			data.Data = append(data.Data, currentDataset)

			// update loop variables
			currentTeam += 1
			currentDataset = scoresDataset{}
			currentDataset.Label = fmt.Sprintf(labelFmt, currentTeam)
			currentDataset.Data = make([]uint64, 0)
			currentPoints = 0
		}

		// add scanned value into current dataset
		currentPoints += points
		currentDataset.Data = append(currentDataset.Data, currentPoints)
	}

	// append final dataset
	data.Data = append(data.Data, currentDataset)

	// convert to the necessary json strings
	labelBytes, err := json.Marshal(data.Labels)
	if err != nil {
		return nil, err
	}
	dataBytes, err := json.Marshal(data.Data)
	if err != nil {
		return nil, err
	}

	// convert to string and return
	retval := new(ScoresData)
	retval.Labels = template.JS(labelBytes)
	retval.Datasets = template.JS(dataBytes)
	return retval, nil
}
