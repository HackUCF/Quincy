package agent

import (
	"encoding/json"
	"net/http"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/sinks/postgres/agent"
	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/common/types"
	"github.com/gin-gonic/gin"
)

func validateScore(types.Score) error {
	return nil
}

// AddScore is a POST route for submitting a completed scorecheck,
// only meant for interaction from an agent.
//
//	@Summary		Submit a completed score check
//	@Description	Accepts a completed score check result from an agent and persists it to the database.
//	@Tags			agent
//	@Accept			json
//	@Produce		json
//	@Param			score	body		types.Score	true	"Completed score check"
//	@Success		200		{object}	object
//	@Failure		400		{object}	object
//	@Router			/agent/completed-score [post]
func AddScore(c *gin.Context) {

	// grab globals from gin context
	sinks := config.Get(c).Sinks
	db, err := conn.GetE(c)

	// check for db errors, fail only if the db sink is enabled.
	// if this fails `db` is safely null and will be ignored by AddScore.
	if err != nil && sinks.DBEnabled() {
		resp := gin.H{
			"message": "failed to get database connection from request context",
			"error":   err,
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// load score from request body
	var score types.Score
	err = json.NewDecoder(c.Request.Body).Decode(&score)
	if err != nil {
		resp := gin.H{
			"message": "couldn't marshall json from request body",
			"error":   err,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// validate score
	if err := validateScore(score); err != nil {
		resp := gin.H{
			"message": "score failed to verify",
			"error":   err,
			"score":   score,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	err = agent.AddScore(c.Request.Context(), db, score)
	if err != nil {
		resp := gin.H{
			"message": "failed to add score to database",
			"error":   err,
			"score":   score,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "score added successfully", "score": score})
}
