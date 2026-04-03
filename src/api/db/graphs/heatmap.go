package graphs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"slices"

	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/types"
)

// heatmapPoint is a single cell in the heatmap chart.
// this is the format the js library expects.
type heatmapPoint struct {
	X string  `json:"x"`
	Y string  `json:"y"`
	V float64 `json:"v"` // uptime percentage in [0, 1]
}

// rawHeatmapData contains the data required to render a heatmap chart.
// it's not rendered yet.
type rawHeatmapData struct {
	XLabels    []string
	YLabels    []string
	DataPoints []heatmapPoint
}

// HeatmapData contains the rendered strings required for the heatmap graph template.
type HeatmapData struct {
	// list of "boxid-serviceid"
	XLabels template.JS
	// list of team labels "team xx"
	YLabels template.JS
	// the datapoints, a list of objects containing x, y, and uptime percentage v.
	DataPoints template.JS
}

// GetHeatmapData returns the template data needed to render a heatmap chart.
// data contains rendered json strings.
// each cell value is the historical uptime percentage for that team/box/service combination.
func GetHeatmapData(db *sql.DB) (*HeatmapData, error) {

	rows, err := db.Query(`
		SELECT service, box, team_num, passed, total
		FROM final_scores
		ORDER BY team_num, box, service
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := new(rawHeatmapData)
	data.XLabels = make([]string, 0)
	data.YLabels = make([]string, 0)
	data.DataPoints = make([]heatmapPoint, 0)

	for rows.Next() {
		// scan the aggregated scores out
		var s types.Score
		var passed, total uint64
		rows.Scan(&s.ServiceName, &s.BoxName, &s.TeamNum, &passed, &total)

		// get x, y, and v
		var point heatmapPoint
		point.X = string(s.BoxName) + "-" + string(s.ServiceName)
		point.Y = fmt.Sprintf(labelFmt, s.TeamNum)
		if total > 0 {
			point.V = float64(passed) / float64(total)
		} else {
			point.V = 0
		}

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
	retval := new(HeatmapData)
	retval.XLabels = template.JS(xBytes)
	retval.YLabels = template.JS(yBytes)
	retval.DataPoints = template.JS(datasetBytes)

	log.Info(string(retval.XLabels))
	log.Info(string(retval.YLabels))
	log.Info(string(retval.DataPoints))

	return retval, nil
}
