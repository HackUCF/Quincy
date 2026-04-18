package cmd

import (
	"github.com/HackUCF/quincy/agent"
	"github.com/spf13/cobra"
)

func agentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Scoring agent",
	}

	cmd.AddCommand(agentStartCmd())

	return cmd
}

func agentStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the scoring agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiURL, _ := cmd.Flags().GetString("api-url")
			loopTime, _ := cmd.Flags().GetInt("loop-time")
			numThreads, _ := cmd.Flags().GetInt("num-threads")

			agent.Start(&agent.AgentConfig{
				APIURL:     apiURL,
				LoopTime:   loopTime,
				NumThreads: numThreads,
			})
			return nil
		},
	}

	cmd.Flags().String("api-url", "http://127.0.0.1:8888", "URL of the API server")
	cmd.Flags().Int("loop-time", 5, "seconds between scoring loops")
	cmd.Flags().Int("num-threads", 10, "number of concurrent scoring threads")

	return cmd
}
