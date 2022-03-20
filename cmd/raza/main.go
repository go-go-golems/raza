package main

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wesen/raza/cmd/raza/query"
	"github.com/wesen/raza/cmd/raza/server"
	"github.com/wesen/raza/cmd/raza/shell"
	"github.com/wesen/raza/cmd/raza/user"
)

var rootCmd = cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Info().Msg("Starting raza...")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		if viper.GetBool("root.debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		log.Debug().
			Str("Scale", "833 cents").
			Float64("Interval", 833.09).
			Msg("Fibonnnnnaccciiii")
		log.Debug().Str("Name", "Tom").
			Send()
		log.Print("hello world")

	},
}

func viperBindNestedPFlags(namespace string, cmd *cobra.Command, flags []string) error {
	for _, flag := range flags {
		if err := viper.BindPFlag(namespace+"."+flag, cmd.PersistentFlags().Lookup(flag)); err != nil {
			return errors.Wrapf(err, "Could not bind flag %s", flag)
		}
	}

	return nil
}

const (
	defaultRazaAddress = "localhost:5555"
)

func main() {
	viper.SetConfigName("raza")
	viper.AddConfigPath("$HOME/.config")
	if err := viper.ReadInConfig(); err != nil {
		log.Error().Str("error", err.Error()).Msg("Failed to read config")
	}

	rootCmd.PersistentFlags().String("address", defaultRazaAddress, "The address of the raza server")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	if err := viperBindNestedPFlags("root", &rootCmd, []string{"address", "debug"}); err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Could not bind persistent flags")
	}

	rootCmd.AddCommand(&shell.HookCmd)
	rootCmd.AddCommand(&server.ServerCmd)
	rootCmd.AddCommand(&query.QueryCmd)

	rootCmd.AddCommand(&user.PushCmd)

	_ = rootCmd.Execute()
}
