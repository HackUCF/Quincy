package testutil

import (
	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/conn"
	"github.com/HackUCF/quincy/api/routes"
	"github.com/HackUCF/quincy/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewTestRouter builds a Gin engine using the real application routes,
// with DB and config injected directly (no globals required).
func NewTestRouter(pool *pgxpool.Pool, cfg *config.APIConfigSpec) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(
		middleware.Recovery(false),
		func(c *gin.Context) {
			c.Set(conn.DBKey, pool)
			c.Set(config.CfgKey, cfg)
			c.Next()
		},
	)
	routes.RegisterRoutes(r)
	return r
}
