package cmd

import (
	"fmt"

	"github.com/HackUCF/Quincy/api"
	"github.com/HackUCF/Quincy/api/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func apiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Scoring API server",
	}

	cmd.PersistentFlags().StringP("config", "c", "config.yaml", "config file path")

	cmd.AddCommand(
		apiStartCmd(),
		apiConfigCmd(),
	)

	return cmd
}

func apiStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgFile, _ := cmd.Flags().GetString("config")
			viper.SetConfigFile(cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("failed to read config file %q: %w", cfgFile, err)
			}

			var cfg config.APIConfigSpec
			if err := viper.Unmarshal(&cfg); err != nil {
				return fmt.Errorf("failed to unmarshal config: %w", err)
			}
			return api.Start(&cfg)
		},
	}
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
