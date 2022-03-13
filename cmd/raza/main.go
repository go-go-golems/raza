package main

import (
	"github.com/spf13/cobra"
	"github.com/wesen/raza/cmd/raza/query"
	"github.com/wesen/raza/cmd/raza/server"
	"github.com/wesen/raza/cmd/raza/shell"
	"github.com/wesen/raza/cmd/raza/user"
)

var rootCmd = cobra.Command{}

const (
	defaultRazaAddress = "localhost:5555"
)

func main() {
	rootCmd.PersistentFlags().String("address", defaultRazaAddress, "The address of the raza server")

	rootCmd.AddCommand(&shell.HookCmd)
	rootCmd.AddCommand(&server.ServerCmd)
	rootCmd.AddCommand(&query.QueryCmd)

	rootCmd.AddCommand(&user.PushCmd)

	_ = rootCmd.Execute()
}
