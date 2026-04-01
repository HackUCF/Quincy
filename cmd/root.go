package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quincy",
		Short: "Cybersecurity competition scoring engine",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetEnvPrefix("QU")
			viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
			viper.AutomaticEnv()
			return nil
		},
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
