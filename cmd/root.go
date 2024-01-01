package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Path string

var rootCmd = &cobra.Command{
	Use:     "sesh",
	Version: "0.0.1",
	Long:    "A nice TUI for managing your TMUX sessions",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Path, "path", "p", "~", "base path for sessions")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
