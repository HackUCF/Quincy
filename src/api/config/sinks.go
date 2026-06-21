package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Sinks) DBEnabled() bool {
	return (PGConfig{}) != s.PGConfig
}

func (s *Sinks) OTelEnabled() bool {
	return s.OTelConfig.Endpoint != ""
}

// DBOr501 returns the original endpoint if the db is configured. Otherwise it returns a simple 501 not implemented.
// This is used to make endpoints that require the db sink to be enabled.
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
// This is used to make middlewares that require the db sink to be enabled.
func (s *Sinks) DBOrNOP(endpoint gin.HandlerFunc) gin.HandlerFunc {

	if s.DBEnabled() {
		return endpoint
	} else {
		return func(c *gin.Context) {}
	}
}
