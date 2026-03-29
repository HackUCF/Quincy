package cmd

import (
	"github.com/HackUCF/Quincy/api"
	"github.com/spf13/cobra"
)

func apiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Scoring API server",
		// no "Run" means print help command
	}

	cmd.AddCommand(
		apiStartCmd(),
		apiConfigCmd(),
	)

	return cmd
}

func apiStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the API server",
		Run: func(cmd *cobra.Command, args []string) {
			api.Start()
		},
	}

	return cmd
}

func apiConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump-config",
		Short: "Dump the default YAML config file",
		Run: func(cmd *cobra.Command, args []string) {
			force, _ := cmd.Flags().GetBool("force")
			print, _ := cmd.Flags().GetBool("print")
			api.DumpConfig(force, print)
		},
	}

	cmd.Flags().BoolP("force", "f", false, "overwrite the config file if it already exists")
	cmd.Flags().BoolP("print", "p", false, "print the config to stdout instead")

	return cmd
}
