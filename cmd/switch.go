package cmd

import (
	"github.com/spf13/cobra"
	filesystem "github.com/xemotrix/sesh/internal/file_system"
	"github.com/xemotrix/sesh/internal/switcher"
)

func init() {
	rootCmd.AddCommand(switchCmd)
}

var switchCmd = &cobra.Command{
	Use:     "switch <path>",
	Aliases: []string{"s", "sw"},
	Example: "sesh switch ~/repos",
	Long:    `switches to an existing tmux session or creates a new one based on a directory under the provided path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filesystem.CleanPath(Path)
		if err != nil {
			return err
		}
		return switcher.InitBubbleTea(path)
	},
}
