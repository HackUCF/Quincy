package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/HackUCF/Quincy/api/config"
)

// DumpConfig writes the default YAML config to disk, or prints it to stdout.
func DumpConfig(forceOverwrite bool, printInstead bool) {
	if printInstead {
		fmt.Println(string(config.DefaultConfigBytes))
		return
	}

	_, err := os.Stat(config.DefaultConfigFile)
	fileExists := !errors.Is(err, os.ErrNotExist)

	if fileExists && !forceOverwrite {
		fmt.Println("The config file already exists. Use --force / -f to overwrite.")
		os.Exit(1)
	}

	err = os.WriteFile(config.DefaultConfigFile, config.DefaultConfigBytes, 0644)
	if err != nil {
		fmt.Printf("Failed to write the config file: %v\n", err)
		os.Exit(1)
	}
}
