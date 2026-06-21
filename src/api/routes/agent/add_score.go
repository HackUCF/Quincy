package agent

import (
	"encoding/json"
	"net/http"

	"github.com/HackUCF/quincy/api/db/agent"
	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/common/types"
	"github.com/gin-gonic/gin"
)

func validateScore(types.Score) error {
	return nil
}

/*
AddScore is a POST route for submitting a completed scorecheck,
only meant for interaction from an agent.
*/
func AddScore(c *gin.Context) {
	db := conn.Get(c)

	// load score from request body
	var score types.Score
	err := json.NewDecoder(c.Request.Body).Decode(&score)
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
