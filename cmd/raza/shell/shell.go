// Package shell implements the low-level commands used to wrap the actual shell history
// for raza. They are expected to be called by shell-specific hooks.
//
// See the zsh implementation in share/init.zsh
package shell

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wesen/raza/cmd/raza/helpers"
	"github.com/wesen/raza/lib/raza"
	"strings"
	"time"
)

var HookCmd = cobra.Command{
	Use:   "_shell",
	Short: "Low-level shell hook function",
}

var AddHookCmd = cobra.Command{
	Use:   "add",
	Short: "Add a new command to the shell history",
	Run: func(cmd *cobra.Command, args []string) {
		session, _ := cmd.Flags().GetInt64("session")
		command, _ := cmd.Flags().GetString("cmd")
		pwd, _ := cmd.Flags().GetString("pwd")
		startTime, _ := cmd.Flags().GetInt64("start_time")

		c, ctx, closeClient, err := helpers.NewRazaShellWrapperClient(cmd)
		if err != nil {
			return
		}
		defer closeClient()

		_, err = c.StartCommand(ctx, &raza.StartCommandRequest{
			SessionId:       session,
			Cmd:             strings.TrimSpace(command),
			Pwd:             pwd,
			StartTimeEpochS: startTime,
			Metadata:        nil,
		})
		if err != nil {
			cmd.PrintErrln("Could not send StartCommand to server", err)
			return
		}
	},
}

var PreHookCmd = cobra.Command{
	Use:   "pre",
	Short: "Run as a pre hook to compute command time and retval",
	Run: func(cmd *cobra.Command, args []string) {
		session, _ := cmd.Flags().GetInt64("session")
		endTime, _ := cmd.Flags().GetInt64("end_time")
		retval, _ := cmd.Flags().GetInt32("retval")

		c, ctx, closeClient, err := helpers.NewRazaShellWrapperClient(cmd)
		if err != nil {
			return
		}
		defer closeClient()

		_, err = c.EndSessionsLastCommand(ctx, &raza.EndSessionsLastCommandRequest{
			SessionId:     session,
			EndTimeEpochS: endTime,
			Retval:        retval,
		})
		if err != nil {
			cmd.PrintErrln("Could not send EndCommand to server", err)
			return
		}
	},
}

var StartSessionCmd = cobra.Command{
	Use:   "start-session",
	Short: "Create a new session and return a session ID",
	Run: func(cmd *cobra.Command, args []string) {
		hostname, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")

		c, ctx, closeClient, err := helpers.NewRazaShellWrapperClient(cmd)
		if err != nil {
			return
		}
		defer closeClient()

		res, err := c.StartSession(ctx, &raza.StartSessionRequest{
			Host: hostname,
			User: user,
		})
		if err != nil {
			cmd.PrintErrln("Could not send StartSession to server", err)
			return
		}
		fmt.Printf("%d\n", res.SessionId)
	},
}

var EndSessionCmd = cobra.Command{
	Use:   "end-session",
	Short: "End a session",
	Run: func(cmd *cobra.Command, args []string) {
		session, _ := cmd.Flags().GetInt64("session")

		c, ctx, closeClient, err := helpers.NewRazaShellWrapperClient(cmd)
		if err != nil {
			return
		}
		defer closeClient()

		_, err = c.EndSession(ctx, &raza.EndSessionRequest{
			SessionId:     session,
			EndTimeEpochS: time.Now().Unix(),
		})
		if err != nil {
			cmd.PrintErrln("Could not send EndSession to server", err)
			return
		}
	},
}

func init() {
	HookCmd.AddCommand(&AddHookCmd)
	HookCmd.AddCommand(&PreHookCmd)
	HookCmd.AddCommand(&StartSessionCmd)
	HookCmd.AddCommand(&EndSessionCmd)

	HookCmd.PersistentFlags().Int64("session", -1, "Session ID")

	AddHookCmd.Flags().String("cmd", "", "Command to add to the shell history")
	_ = AddHookCmd.MarkFlagRequired("cmd")
	AddHookCmd.Flags().String("pwd", "", "Current working directory")
	_ = AddHookCmd.MarkFlagRequired("pwd")
	AddHookCmd.Flags().Int64("start_time", 0, "Start time of the command (epoch seconds)")
	_ = AddHookCmd.MarkFlagRequired("start_time")

	PreHookCmd.Flags().Int64("end_time", 0, "End time of the command (epoch seconds)")
	_ = PreHookCmd.MarkFlagRequired("end_time")
	PreHookCmd.Flags().Int32("retval", 0, "Return value of the command")
	_ = PreHookCmd.MarkFlagRequired("retval")

	StartSessionCmd.Flags().String("host", "", "Hostname of the machine")
	_ = StartSessionCmd.MarkFlagRequired("host")
	StartSessionCmd.Flags().String("user", "", "User running the command")
	_ = StartSessionCmd.MarkFlagRequired("user")
}
