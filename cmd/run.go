package cmd

import (
	"github.com/edup2p/relay-server/config"
	"github.com/edup2p/relay-server/relay"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the relay server",
	// Long: ``,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.ReadConfig(cfgFile)
		if err != nil {
			return err
		}

		privKey, err := config.ReadKey(cfg.KeyFile)
		if err != nil {
			return err
		}

		return relay.RunRelay(*cfg, *privKey)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
