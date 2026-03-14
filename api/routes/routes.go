/*
Package routes contains all of the Gin HTTP API logic.
This is includes which functions are responsible for which paths and methods,
user input parsing and validation,
as well as request logging and error handling.
Most Gin handler functions are located in subpackages.
*/
package routes

import (
	"github.com/HackUCF/Quincy/api/routes/graphs"
	"github.com/HackUCF/Quincy/api/routes/misc"
	"github.com/HackUCF/Quincy/api/routes/scoring"
	"github.com/HackUCF/Quincy/api/routes/users"
	"github.com/HackUCF/Quincy/common/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// create a gin router
func initRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.Recovery)
	router.Use(middleware.Logging)
	router.Use(cors.Default()) // insecure, allows all cross-origin requests

	v1 := router.Group("/api/v1")
	{
		scoringGroup := v1.Group("/scores")
		{
			scoringGroup.POST("/", scoring.AddScore)                 // /api/v1/scores
			scoringGroup.GET("/team", scoring.GetTeamScores)         // /api/v1/scores/team
			scoringGroup.GET("/box", scoring.GetBoxScores)           // /api/v1/scores/box
			scoringGroup.GET("/service", scoring.GetServiceScores)   // /api/v1/scores/service
			scoringGroup.GET("/current", scoring.GetRecentChecks)    // /api/v1/scores/current
			scoringGroup.GET("/detailed", scoring.GetDetailedScores) // /api/v1/scores/detailed
		}

		userGroup := v1.Group("/users")
		{
			userGroup.GET("", users.GetAllUsers) // /api/v1/users
			userGroup.POST("", users.SubmitPCR)  // /api/v1/users
		}

		checkGroup := v1.Group("/checks")
		{
			checkGroup.GET("", scoring.GetCheck) // /api/v1/checks
		}

		graphsGroup := v1.Group("/graphs")
		{
			graphsGroup.GET("scoreboard", graphs.GetScoreboard) // /api/v1/graphs/scoreboard
			graphsGroup.GET("scores", graphs.GetScores)         // /api/v1/graphs/scores
			graphsGroup.GET("standings", graphs.GetStandings)   // /api/v1/graphs/standings
			graphsGroup.GET("heatmap", graphs.GetHeatmap)       // /api/v1/graphs/heatmap
		}

		v1.GET("/config", misc.GetConfig) // /api/v1/config
	}

	router.NoRoute(misc.NoRoute(router))

	return router
}
