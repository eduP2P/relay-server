package cmd

import (
	"log/slog"
	"os"
	"strings"

	"github.com/edup2p/common/types"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "edup2p_relay_server",
	Short: "EduP2P/ToverSok relay server",
	// Long: ``,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		switch logLevel {
		case "trace":
			slog.SetLogLoggerLevel(types.LevelTrace)
		case "debug":
			slog.SetLogLoggerLevel(slog.LevelDebug)
		case "info":
			slog.SetLogLoggerLevel(slog.LevelInfo)
		case "warn":
			slog.SetLogLoggerLevel(slog.LevelWarn)
		case "error":
			slog.SetLogLoggerLevel(slog.LevelError)
		default:
			slog.Warn("could not parse log level (try: trace, debug, info, warn, error)", "level", logLevel)
			return
		}
		println("Logging level set to", strings.ToUpper(logLevel))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	cfgFile  string
	logLevel string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "relay.toml", "config file path")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level")
}
