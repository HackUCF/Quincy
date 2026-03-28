package cmd

import (
	"github.com/HackUCF/Quincy/agent"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Scoring agent",
}

var agentStartCmd = &cobra.Command{
	Run:   agent.Start,
	Use:   "start",
	Short: "Start the scoring agent",
}

func initAgentCmd() {
	agentCmd.AddCommand(agentStartCmd)
}
