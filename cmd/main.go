package cmd

import (
	"os"
)

func init() {
	initAgentCmd()
	initApiCmd()
	initRoot()

	rootCmd.AddCommand(agentCmd, apiCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
