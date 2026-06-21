/*
Package routes contains all of the Gin HTTP API logic.
This is includes which functions are responsible for which paths and methods,
user input parsing and validation,
as well as request logging and error handling.
Most Gin handler functions are located in subpackages.
*/
package routes

import (
	"github.com/HackUCF/quincy/api/config"
	_ "github.com/HackUCF/quincy/api/openapi"
	"github.com/HackUCF/quincy/api/routes/agent"
	"github.com/HackUCF/quincy/api/routes/graphs"
	"github.com/HackUCF/quincy/api/routes/misc"
	"github.com/HackUCF/quincy/api/routes/scoring"
	"github.com/HackUCF/quincy/api/routes/users"
	"github.com/HackUCF/quincy/api/sinks/postgres/conn"
	"github.com/HackUCF/quincy/common/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes registers all API routes on router. Called by initRoutes and tests.
// Routes that require the db are replaced with 501 not implemented.
func RegisterRoutes(router *gin.Engine, s config.Sinks) {
	v1 := router.Group("/api/v1")
	{
		agentGroup := v1.Group("/agent")
		{
			agentGroup.GET("/new-check", agent.GetCheck)        // /api/v1/agent/new-check
			agentGroup.POST("/completed-score", agent.AddScore) // /api/v1/agent/completed-score
		}

		scoringGroup := v1.Group("/scores")
		{
			scoringGroup.GET("/team", s.DBOr501(scoring.GetTeamScores))         // /api/v1/scores/team
			scoringGroup.GET("/box", s.DBOr501(scoring.GetBoxScores))           // /api/v1/scores/box
			scoringGroup.GET("/service", s.DBOr501(scoring.GetServiceScores))   // /api/v1/scores/service
			scoringGroup.GET("/current", s.DBOr501(scoring.GetRecentChecks))    // /api/v1/scores/current
			scoringGroup.GET("/detailed", s.DBOr501(scoring.GetDetailedScores)) // /api/v1/scores/detailed
		}

		userGroup := v1.Group("/users")
		{
			userGroup.GET("", users.GetAllUsers)           // /api/v1/users
			userGroup.POST("", s.DBOr501(users.SubmitPCR)) // /api/v1/users
		}

		graphsGroup := v1.Group("/graphs")
		{
			graphsGroup.GET("scoreboard", s.DBOr501(graphs.GetScoreboard)) // /api/v1/graphs/scoreboard
			graphsGroup.GET("scores", s.DBOr501(graphs.GetScores))         // /api/v1/graphs/scores
			graphsGroup.GET("standings", s.DBOr501(graphs.GetStandings))   // /api/v1/graphs/standings
			graphsGroup.GET("heatmap", s.DBOr501(graphs.GetHeatmap))       // /api/v1/graphs/heatmap
		}

		v1.GET("/config", misc.GetConfig) // /api/v1/config
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(misc.NoRoute(router))
}

func initRoutes(sinks config.Sinks) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		middleware.Recovery(false),
		middleware.Logging(),
		middleware.InsecureCORS(),
		config.ConfigMiddleware(),
		sinks.DBOrNOP(conn.DBMiddleware()),
	)
	RegisterRoutes(router, sinks)
	return router
}
