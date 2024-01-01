package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xemotrix/sesh/internal/killer"
)

func init() {
	rootCmd.AddCommand(killCmd)
}

var killCmd = &cobra.Command{
	Use:     "kill",
	Aliases: []string{"k"},
	Example: "sesh kill",
	Long:    "Interactively kill one or many sessions",
	RunE: func(cmd *cobra.Command, args []string) error {

		return killer.InitBubbleTea()
	},
}
