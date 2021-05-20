package cmd

import (
	"github.com/spf13/cobra"

	"gitee.com/xchain/go-xchain/client"
)

// Commands registers a sub-tree of commands to interact with
// local private key storage.
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "start lite node",
		Long:  `start lite node `,
	}
	cmd.AddCommand(
		StartLiteNodeCmd(),
		client.LineBreak,
	)
	return cmd
}
