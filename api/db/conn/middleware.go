package conn

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

// DBKey is the key used to store the database object in the request context.
const DBKey string = "db-object"

// Middleware stores the database connection object in the gin context.
// This can be retrieved as necessary by routes.
func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DBKey, db)
	}
}

// Get returns the database object stored in the request context.
// Panics on error.
func Get(c *gin.Context) *sql.DB {
	return c.MustGet(DBKey).(*sql.DB)
}

// GetE returns the database object stored in the request context.
func GetE(c *gin.Context) (*sql.DB, error) {
	// pull from context
	obj, ok := c.Get(DBKey)
	if !ok {
		err := fmt.Errorf("failed to get database object with key %q", DBKey)
		return nil, err
	}

	// cast as database object
	db, ok := obj.(*sql.DB)
	if !ok {
		err := fmt.Errorf("object with key %q is not a sql database", DBKey)
		return nil, err
	}

	return db, nil
}
