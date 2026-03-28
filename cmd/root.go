package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "quincy",
	Short: "Cybersecurity competition scoring engine",
}

func initRoot() {}
