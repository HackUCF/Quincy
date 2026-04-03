package config

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// CfgKey is the key used to store the config object in the request context.
const CfgKey string = "cfg-object"

// Middleware stores the config connection object in the gin context.
// This can be retrieved as necessary in routes.
func ConfigMiddleware() gin.HandlerFunc {

	// runs at router initialization-time
	if !cfgIsValid.Load() {
		// panic if config not loaded
		log.Panic("config not stored, cannot create middleware")
	}

	// runs during each request
	return func(c *gin.Context) {
		c.Set(CfgKey, cfg)
	}
}

// Get returns the config object stored in the request context.
// Panics on error.
func Get(c *gin.Context) *APIConfigSpec {
	return c.MustGet(CfgKey).(*APIConfigSpec)
}

// GetE returns the config object stored in the request context.
func GetE(c *gin.Context) (*APIConfigSpec, error) {

	// pull from context
	obj, ok := c.Get(CfgKey)
	if !ok {
		err := fmt.Errorf("failed to get config object with key %q", CfgKey)
		return nil, err
	}

	// cast as config object
	cfg, ok := obj.(*APIConfigSpec)
	if !ok {
		err := fmt.Errorf("object with key %q is not a an api configuration", CfgKey)
		return nil, err
	}

	return cfg, nil
}
