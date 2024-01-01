package main

import (
	"fmt"
	"sesh/creator"
	"sesh/switcher"

	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:     "sesh <command> <path>",
		Example: "sesh switch ~/repos\nsesh create ~/repos",
		Version: "0.0.1",
		Long:    "A nice TUI for managing your TMUX sessions",
		ValidArgs: []string{
			"switch",
			"create",
		},
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			command := args[0]
			path := args[1]
			switch command {
			case "switch":
				return switcher.InitBubbleTea(path)
			case "create":
				return creator.InitBubbleTea(path)
			default:
				return fmt.Errorf("Invalid command: %s", command)
			}
		},
	}
	cmd.Execute()
}
