package query

import (
	"github.com/spf13/cobra"
	"github.com/wesen/raza/cmd/raza/helpers"
	"github.com/wesen/raza/lib/raza"
	"io"
)

var QueryCmd = cobra.Command{
	Use:   "query",
	Short: "Query history",
}

var ListSessionsCmd = cobra.Command{
	Use:   "list-sessions",
	Short: "List sessions",
	Run: func(cmd *cobra.Command, args []string) {
		c, ctx, closeClient, err := helpers.NewRazaQueryClient(cmd)
		if err != nil {
			return
		}
		defer closeClient()

		res, err := c.GetSessions(ctx, &raza.GetSessionsRequest{})
		if err != nil {
			cmd.PrintErrf("Could not list sessions: %s\n", err.Error())
		}

		for {
			session, err := res.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				cmd.PrintErrf("Could not list sessions: %s\n", err.Error())
				return
			}
			cmd.Printf("%v\n", session)
		}
	},
}

var ListCommandsCmd = cobra.Command{
	Use:   "list-commands",
	Short: "List commands",
	Run: func(cmd *cobra.Command, args []string) {
		c, ctx, closeClient, err := helpers.NewRazaQueryClient(cmd)
		if err != nil {
			return
		}
		defer closeClient()

		res, err := c.GetCommands(ctx, &raza.GetCommandsRequest{})
		if err != nil {
			cmd.PrintErrf("Could not list sessions\n")
		}

		for {
			command, err := res.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				cmd.PrintErrln("Could not list commands")
				return
			}
			cmd.Printf("%v\n", command)
		}
	},
}

func init() {
	QueryCmd.AddCommand(&ListSessionsCmd)
	QueryCmd.AddCommand(&ListCommandsCmd)
}
