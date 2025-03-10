package cmd

import (
	"log/slog"
	"os"

	"github.com/edup2p/common/types/key"
	"github.com/edup2p/relay-server/config"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Config and key generation",
	Long:  `This command will generate a default config and key.`,
	Run: func(_ *cobra.Command, _ []string) {
		var cfg *config.Config

		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			slog.Info("config file does not exist, creating default config", "file", cfgFile)

			cfg = new(config.Config)
			defaults.SetDefaults(cfg)

			if err := config.WriteConfig(cfg, cfgFile); err != nil {
				slog.Error("error writing config file", "err", err)
				os.Exit(1)
			}

			slog.Info("written default config file", "to", cfgFile)
		} else {
			slog.Error("config file already exists, not overwriting...", "file", cfgFile)

			cfg, err = config.ReadConfig(cfgFile)
			if err != nil {
				slog.Error("error reading config file", "err", err)
				os.Exit(1)
			}
		}

		if _, err := os.Stat(cfg.KeyFile); os.IsNotExist(err) {
			slog.Info("key file does not exist, generating new key...", "keyFile", cfg.KeyFile)
			k := key.NewNode()
			if err := config.WriteKey(k, cfg.KeyFile); err != nil {
				slog.Error("error writing new key file", "err", err)
			}
		} else {
			slog.Warn("key file already exists, skipping overwriting", "keyFile", cfg.KeyFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
