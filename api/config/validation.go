package config

// validate the api configuration file
// tries to prevent panics later down the line
// definitely could be cleaner

import (
	"errors"
	"fmt"
	"net"
	"slices"

	"github.com/HackUCF/Quincy/common/types"
)

// stored in two places: here and in the db schema
const MaxStringLength int = 16

func (cfg *APIConfigSpec) validate() error {
	// low hanging fruit
	if cfg.NumTeams == 0 {
		return errors.New("'num_teams' cannot be zero")
	}

	if cfg.DBFile == "" {
		return errors.New("'db_file' missing or empty")
	}

	if len(cfg.Boxes) == 0 {
		return errors.New("'boxes' missing or empty")
	}

	// create slice to store box ids and display names
	boxIDs := make([]types.BoxID, 0, len(cfg.Boxes))
	boxDNs := make([]string, 0, len(cfg.Boxes))

	// validate each box and add the info to slices
	for _, box := range cfg.Boxes {
		err := box.validate()
		if err != nil {
			return fmt.Errorf("box %q failed to validate: %w", box.ID, err)
		}

		if slices.Contains(boxIDs, box.ID) {
			return fmt.Errorf("two boxes share the same id: %q", box.ID)
		}

		if slices.Contains(boxDNs, box.DisplayName) {
			// should just be a warning
			// two boxes have the same display name
			// could be intentional but would probably just be confusing
		}

		boxIDs = append(boxIDs, box.ID)
		boxDNs = append(boxDNs, box.DisplayName)
	}

	// create slice to store userlist ids and display names
	ulIDs := make([]types.UserListID, 0, len(cfg.UserLists))
	ulDNs := make([]string, 0, len(cfg.UserLists))

	// validate each userlist and add the info to slices
	for _, ul := range cfg.UserLists {
		err := ul.validate()
		if err != nil {
			return fmt.Errorf("user list %q failed to validate: %w", ul.ID, err)
		}

		if slices.Contains(ulIDs, ul.ID) {
			return fmt.Errorf("two userlists share the same id: %q", ul.ID)
		}

		if slices.Contains(ulDNs, ul.DisplayName) {
			// should just be a warning
			// two userlists have the same display name
			// could be intentional but would probably just be confusing
		}

		ulIDs = append(ulIDs, ul.ID)
		ulDNs = append(ulDNs, ul.DisplayName)
	}

	err := cfg.HTTP.validate()
	if err != nil {
		return fmt.Errorf("http config failed to validate: %w", err)
	}

	return nil
}

func (box *BoxSpec) validate() error {
	if box.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if box.DisplayName == "" {
		// should just be a warning. this isn't required for functionality
	}

	if len(box.ID) > MaxStringLength {
		return fmt.Errorf("id is greater than %d characters long", MaxStringLength)
	}

	if len(box.DisplayName) > MaxStringLength {
		return fmt.Errorf("display name %q is greater than %d characters long", box.DisplayName, MaxStringLength)
	}

	if len(box.Services) == 0 {
		return fmt.Errorf("a box cannot have zero services")
	}

	if box.Host == "" {
		// should just be a warning. this isn't required for functionality
	}

	// create slice to store userlist ids and display names
	svcIDs := make([]types.ServiceID, 0, len(box.Services))
	svcDNs := make([]string, 0, len(box.Services))

	// validate each userlist and add the info to slices
	for _, svc := range box.Services {
		err := svc.validate()
		if err != nil {
			return fmt.Errorf("service %q failed to validate: %w", svc.ID, err)
		}

		if slices.Contains(svcIDs, svc.ID) {
			return fmt.Errorf("two services share the same id: %q", svc.ID)
		}

		if slices.Contains(svcDNs, svc.DisplayName) {
			// should just be a warning
			// two services have the same display name
			// could be intentional but would probably just be confusing
		}

		svcIDs = append(svcIDs, svc.ID)
		svcDNs = append(svcDNs, svc.DisplayName)
	}

	return nil
}

func (svc *ServiceSpec) validate() error {
	if svc.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if svc.DisplayName == "" {
		// should just be a warning. this isn't required for functionality
	}

	if svc.CheckID == "" {
		return fmt.Errorf("check id cannot be empty")
	}

	if len(svc.ID) > MaxStringLength {
		return fmt.Errorf("id is greater than %d characters long", MaxStringLength)
	}

	if len(svc.DisplayName) > MaxStringLength {
		return fmt.Errorf("display name %q is greater than %d characters long", svc.DisplayName, MaxStringLength)
	}

	return nil
}

// ensure some length constraints are met
func (ul *UserListSpec) validate() error {
	if ul.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}

	if ul.DisplayName == "" {
		// should just be a warning. this isn't required for functionality
	}

	if len(ul.ID) > MaxStringLength {
		return fmt.Errorf("id is greater than %d characters long", MaxStringLength)
	}

	if len(ul.DisplayName) > MaxStringLength {
		return fmt.Errorf("display name %q is greater than %d characters long", ul.DisplayName, MaxStringLength)
	}

	if len(ul.Users) == 0 {
		return fmt.Errorf("a userlist cannot have zero users")
	}

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
