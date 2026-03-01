package graphs

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/common/types"
)

// StandingsData contains the rendered strings required for the standings graph template.
type StandingsData struct {
	// list of team labels "team xx"
	Labels template.JS
	// checks passed per team, parallel to Labels
	Data template.JS
}

// GetStandingsData returns the template data needed to render a current standings bar chart.
// data contains rendered json strings.
func GetStandingsData() (*StandingsData, error) {
	db := conn.Get()

	// query each teams scores
	rows, err := db.Query(`
		SELECT team_num, SUM(passed)
		FROM final_scores
		GROUP BY team_num
		ORDER BY team_num
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	labels := make([]string, 0)
	data := make([]uint64, 0)

	// read each teams final scores
	for rows.Next() {
		var team types.TeamNum
		var passed uint64

		err := rows.Scan(&team, &passed)
		if err != nil {
			return nil, err
		}

		labels = append(labels, fmt.Sprintf(labelFmt, team))
		data = append(data, passed)
	}

	// render json
	labelBytes, err := json.Marshal(labels)
	if err != nil {
		return nil, err
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// return json strings
	retval := new(StandingsData)
	retval.Labels = template.JS(labelBytes)
	retval.Data = template.JS(dataBytes)
	return retval, nil
}
