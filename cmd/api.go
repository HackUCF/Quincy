package cmd

import (
	"github.com/HackUCF/Quincy/api"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Scoring API server",
}

var apiStartCmd = &cobra.Command{
	Run:   api.Start,
	Use:   "start",
	Short: "Start the API server",
}

var apiConfigCmd = &cobra.Command{
	Run:   api.DumpConfig,
	Use:   "config",
	Short: "Dump the default YAML config file",
}

func initApiCmd() {
	apiConfigCmd.Flags().BoolP("force", "f", false, "overwrite the config file if it already exists")
	apiConfigCmd.Flags().BoolP("print", "p", false, "print the config to stdout instead")

	apiCmd.AddCommand(apiStartCmd)
	apiCmd.AddCommand(apiConfigCmd)
}
