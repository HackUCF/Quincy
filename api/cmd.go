package main

import (
	_ "embed"
	"errors"
	"flag"
	"os"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/common/log"
)

var (
	// store the default config file for easier drag and drop
	//go:embed config.yaml
	defaultConfigData []byte

	// generate the default config file
	dumpConfigLong  = flag.Bool("new-config", false, "dump the default config file to ./config.yaml")
	dumpConfigShort = flag.Bool("c", false, "short version of -new-config")

	// // generate a .env file with default values
	// dumpDotEnvLong  = flag.Bool("new-env", false, "dump an example .env file")
	// dumpDotEnvShort = flag.Bool("e", false, "short version of -new-env")
)

// perform the cli actions.
// handles everything that isn't running the api server.
// panics on error.
func doCLI() {
	flag.Parse()

	dumpConfig := *dumpConfigLong || *dumpConfigShort
	// dumpDotEnv := *dumpDotEnvLong || *dumpDotEnvShort

	// dump the config if requested
	if dumpConfig {

		// panic if file exists
		_, err := os.Stat(config.DefaultConfigFile)
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic(
				"cannot dump config, file already exists",
				"config_file", config.DefaultConfigFile,
			)
		}

		// write to file
		err = os.WriteFile(config.DefaultConfigFile, defaultConfigData, 0644)
		if err != nil {
			log.Panic(
				"failed to write default config file",
				"config_file", config.DefaultConfigFile,
				"error", err,
			)
		}

		// exit program
		log.Info("successfully wrote config file, exiting...")
		os.Exit(0)
	}

	// if dumpDotEnv {
	// 	// ...
	// }
}
