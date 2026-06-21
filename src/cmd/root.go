package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version contains the version number of the build as a string.
// In the repo this is a magic development value.
// CI changes this value during build time.
var Version = "v0.0.0 (development)"

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quincy",
		Short: "Cybersecurity competition scoring engine",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			viper.SetEnvPrefix("QU")
			viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
			viper.AutomaticEnv()
			return nil
		},
	}

	// print the version if requested
	// otherwise just default to help message
	//   ./quincy -v   ->   v0.0.123-abc1234
	//   ./quincy      ->   *help message*
	cmd.Flags().BoolP("version", "v", false, "print version and exit")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		v, _ := cmd.Flags().GetBool("version")
		if v {
			fmt.Println(Version)
			return
		}
		cmd.Help()
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
