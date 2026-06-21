package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Sinks) DBEnabled() bool {
	return (PGConfig{}) != s.PGConfig
}

// DBOr501 returns the original endpoint if the db is configured. Otherwise it returns a simple 501 not implemented.
// This is used to make db 
func (s *Sinks) DBOr501(endpoint gin.HandlerFunc) gin.HandlerFunc {

	if s.DBEnabled() {
		return endpoint
	} else {
		return func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{})
		}
	}
}

// DBOrNOP returns the original endpoint if the db is configured. Otherwise it is a no-op.
// This is used to make db specific middlewares.
func (s *Sinks) DBOrNOP(endpoint gin.HandlerFunc) gin.HandlerFunc {

	if s.DBEnabled() {
		return endpoint
	} else {
		return func(c *gin.Context) {}
	}
}
