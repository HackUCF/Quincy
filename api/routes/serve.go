package routes

import (
	"fmt"
	"net"
	"strconv"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/common/log"
)

// ServeRoutes generates an HTTP router and starts listening!
// Requires the config to be loaded.
func ServeRoutes() error {
	router := initRoutes()

	if router == nil {
		return fmt.Errorf("cannot serve from nil router")
	}

	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	conn_string := net.JoinHostPort(cfg.HTTP.Host, strconv.Itoa(cfg.HTTP.Port))

	log.Info(
		"starting to serve http",
		"host", cfg.HTTP.Host,
		"port", cfg.HTTP.Port,
		"conn_string", conn_string,
	)

	return router.Run(conn_string)
}
