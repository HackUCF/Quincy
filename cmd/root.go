package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quincy",
		Short: "Cybersecurity competition scoring engine",
		// no "Run" means print help command
	}
	
	cmd.AddCommand(
		apiCmd(),
		agentCmd(),
	)

	return cmd
}

func Execute() {
	err := rootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
