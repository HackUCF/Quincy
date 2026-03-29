package cmd

import (
	"github.com/HackUCF/Quincy/agent"
	"github.com/spf13/cobra"
)

func agentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Scoring agent",
		// no "Run" means print help command
	}

	cmd.AddCommand(agentStartCmd())

	return cmd
}

func agentStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the scoring agent",
		Run: func(cmd *cobra.Command, args []string) {
			agent.Start()
		},
	}

	return cmd
}
