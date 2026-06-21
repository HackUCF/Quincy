package conn

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBKey is the key used to store the database object in the request context.
const DBKey string = "db-object"

// Middleware stores the database connection object in the gin context.
// This can be retrieved as necessary by routes.
func DBMiddleware(c *gin.Context) {

	c.Set(DBKey, db)
}

// Get returns the database object stored in the request context.
// Panics on error.
func Get(c *gin.Context) *pgxpool.Pool {

	return c.MustGet(DBKey).(*pgxpool.Pool)
}

// GetE returns the database object stored in the request context.
func GetE(c *gin.Context) (*pgxpool.Pool, error) {

	// pull from context
	obj, ok := c.Get(DBKey)
	if !ok {
		err := fmt.Errorf("failed to get database object with key %q", DBKey)
		return nil, err
	}

	// cast as database object
	db, ok := obj.(*pgxpool.Pool)
	if !ok {
		err := fmt.Errorf("object with key %q is not a sql database", DBKey)
		return nil, err
	}

	return db, nil
}
