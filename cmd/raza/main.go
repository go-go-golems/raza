package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var shellHookCmd = cobra.Command{
	Use:   "_shell",
	Short: "Low-level shell hook function",
}

var shellAddHookCmd = cobra.Command{
	Use:   "add",
	Short: "Add a new command to the shell history",
	Run: func(cmd *cobra.Command, args []string) {
		session, _ := cmd.Flags().GetString("session")
		command, _ := cmd.Flags().GetString("cmd")
		pwd, _ := cmd.Flags().GetString("pwd")
		startTime, _ := cmd.Flags().GetInt("start_time")
		fmt.Printf("cmd: '%s', session: %s, pwd: %s, start_time: %d\n", strings.TrimSpace(command), session, pwd, startTime)
	},
}

var shellPreHookCmd = cobra.Command{
	Use:   "pre",
	Short: "Run as a pre hook to compute command timem and retval",
	Run: func(cmd *cobra.Command, args []string) {
		session, _ := cmd.Flags().GetString("session")
		endTime, _ := cmd.Flags().GetInt("end_time")
		retval, _ := cmd.Flags().GetInt("retval")
		fmt.Printf("session: %s, end_time: %d, retval: %d\n", session, endTime, retval)
	},
}

var shellStartSessionCmd = cobra.Command{
	Use:   "start-session",
	Short: "Create a new session and return a session ID",
	Run: func(cmd *cobra.Command, args []string) {
		hostname, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		fmt.Printf("_raza_session_%s_%s", hostname, user)
	},
}

var rootCmd = cobra.Command{}

func main() {
	rootCmd.AddCommand(&shellHookCmd)
	shellHookCmd.AddCommand(&shellAddHookCmd)
	shellHookCmd.AddCommand(&shellPreHookCmd)
	shellHookCmd.AddCommand(&shellStartSessionCmd)

	shellHookCmd.PersistentFlags().String("session", "", "Session ID")

	shellAddHookCmd.Flags().String("cmd", "", "Command to add to the shell history")
	_ = shellAddHookCmd.MarkFlagRequired("cmd")
	shellAddHookCmd.Flags().String("pwd", "", "Current working directory")
	_ = shellAddHookCmd.MarkFlagRequired("pwd")
	shellAddHookCmd.Flags().Int("start_time", 0, "Start time of the command (epoch seconds)")
	_ = shellAddHookCmd.MarkFlagRequired("start_time")

	shellPreHookCmd.Flags().Int("end_time", 0, "End time of the command (epoch seconds)")
	_ = shellPreHookCmd.MarkFlagRequired("end_time")
	shellPreHookCmd.Flags().Int("retval", 0, "Return value of the command")
	_ = shellPreHookCmd.MarkFlagRequired("retval")

	shellStartSessionCmd.Flags().String("host", "", "Hostname of the machine")
	_ = shellStartSessionCmd.MarkFlagRequired("host")
	shellStartSessionCmd.Flags().String("user", "", "User running the command")
	_ = shellStartSessionCmd.MarkFlagRequired("user")

	_ = rootCmd.Execute()
}
