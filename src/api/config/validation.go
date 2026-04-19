package config

// validate the api configuration file
// tries to prevent panics later down the line
// definitely could be cleaner

import (
	"fmt"
	"net"
)

// stored in two places: here and in the db schema
const MaxStringLength int = 16

// Validate validates the API configuration.
func (cfg *APIConfigSpec) Validate() error {
	return nil
}

// TODO
func (box *BoxSpec) validate() error {
	return nil
}

// TODO
func (ul *UserListSpec) validate() error {
	return nil
}

func (http *HTTPSpec) validate() error {
	if http.Port < 0 || http.Port > 65535 {
		return fmt.Errorf("tcp port %d is invalid", http.Port)
	}

	// probably get rid of this?
	// listening on a hostname is normal i think
	ip := net.ParseIP(http.Host)
	if ip == nil {
		return fmt.Errorf("listening host %q is invalid, should parse as an ip address", http.Host)
	}

	return nil
}
