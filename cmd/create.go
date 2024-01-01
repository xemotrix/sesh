package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xemotrix/sesh/internal/creator"
	filesystem "github.com/xemotrix/sesh/internal/file_system"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create",
	Example: "sesh create -p ~/repos",
	Long:    "create a new directory+session under the provided path",
	Aliases: []string{"c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filesystem.CleanPath(Path)
		if err != nil {
			return err
		}
		return creator.InitBubbleTea(path)
	},
}
