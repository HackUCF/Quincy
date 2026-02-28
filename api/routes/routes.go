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
	"github.com/HackUCF/Quincy/common/log"
	"github.com/HackUCF/Quincy/common/middleware"
	"github.com/gin-gonic/gin"
)

// create a gin router
func initRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.Recovery)
	router.Use(middleware.Logging)

	router.NoRoute(misc.NoRoute)

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
			// put a check request
		}

		graphsGroup := v1.Group("/graphs")
		{
			graphsGroup.GET("scoreboard", graphs.GetScoreboard)
		}

		v1.GET("/config", misc.GetConfig) // /api/v1/config
	}

	// lol
	router.GET("/panic", func(c *gin.Context) {
		log.Debug("debug log test")
		log.Info("info log test")
		log.Warn("warn log test")
		log.Error("error log test")
		log.Panic("user requested forced panic endpoint")
	})

	return router
}
