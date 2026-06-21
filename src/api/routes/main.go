package routes

import (
	"fmt"
	"net"
	"strconv"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/log"
)

// ServeRoutes generates an HTTP router and starts listening!
// Requires the config to be loaded.
func ServeRoutes(cfg *config.APIConfigSpec) error {
	router := initRoutes(cfg.Sinks)

	if router == nil {
		return fmt.Errorf("cannot serve from nil router")
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
