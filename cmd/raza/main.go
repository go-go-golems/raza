package main

import (
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wesen/raza/cmd/raza/query"
	"github.com/wesen/raza/cmd/raza/server"
	"github.com/wesen/raza/cmd/raza/shell"
	"github.com/wesen/raza/cmd/raza/user"
	"os"
	"time"
)

var rootCmd = cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		if viper.GetBool("root.debug") {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		file := viper.GetString("root.log-file")
		if file == "" {
			if isatty.IsTerminal(os.Stderr.Fd()) {
				log.Debug().Msg("stderr is a terminal")
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
			} else {
				log.Debug().Msg("stderr is not a terminal")
				log.Logger = log.Output(os.Stderr)
			}
		} else {
			w, err := os.Open(file)
			if err != nil {
				log.Fatal().Err(err).Msgf("Could not open log file %s", file)
			}
			log.Debug().Str("log-file", file).Msg("Logging to file")
			log.Logger = log.Output(w)
		}

		if viper.GetBool("root.log-line") {
			log.Logger = log.With().Caller().Logger()
		}

		if viper.GetBool("root.log-error-stacktrace") {
			log.Debug().Msg("Logging error stacktraces")
			zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		}
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
		log.Debug().Err(err).Msg("Failed to read config")
	}

	rootCmd.PersistentFlags().String("address", defaultRazaAddress, "The address of the raza server")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")
	rootCmd.PersistentFlags().Bool("log-error-stacktrace", false, "Enable stacktrace logging on errors")
	rootCmd.PersistentFlags().Bool("log-line", true, "Enable logging of file ane line number")
	rootCmd.PersistentFlags().String("log-file", "", "Enable logging to file")
	if err := viperBindNestedPFlags("root", &rootCmd,
		[]string{"address", "debug", "log-error-stacktrace", "log-line", "log-file"}); err != nil {
		log.Fatal().Err(err).Msg("Could not bind persistent flags")
	}

	rootCmd.AddCommand(&shell.HookCmd)
	rootCmd.AddCommand(&server.ServerCmd)
	rootCmd.AddCommand(&query.QueryCmd)

	rootCmd.AddCommand(&user.PushCmd)

	_ = rootCmd.Execute()
}
