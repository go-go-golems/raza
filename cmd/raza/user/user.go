// Package user contains user facing helper commands.
// These are most commands in order to avoid too much hierarchy
package user

import "github.com/spf13/cobra"

var PushCmd = cobra.Command{
	Use:   "push",
	Short: "Push metadata on the stack",
	Long:  `Push a list of key=value pairs onto the metadata stack. These will be added to the metadata sent out for each command, until the metadata is popped.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}
