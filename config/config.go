package config

import (
	"fmt"
	"io"
	"net/netip"
	"os"

	"github.com/edup2p/common/types/key"
	"github.com/mcuadros/go-defaults"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Bind     netip.Addr `toml:"bind"      default:"0.0.0.0"   comment:"server bind address for both HTTP and STUN"`
	Port     uint16     `toml:"port"      default:"80"        comment:"port for HTTP"`
	STUNPort uint16     `toml:"stun_port" default:"3478"      comment:"port for STUN"`
	KeyFile  string     `toml:"key_file"  default:"relay.key" comment:"path to private key file"`
}

func ReadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := new(Config)
	defaults.SetDefaults(cfg)
	err = toml.Unmarshal(b, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return cfg, nil
}

func WriteConfig(config *Config, configFile string) error {
	b, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	file, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(b); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func ReadKey(keyFile string) (*key.NodePrivate, error) {
	file, err := os.Open(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open key file: %w", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	if len(b) != key.Len {
		return nil, fmt.Errorf("invalid file byte key length (%d != %d)", len(b), key.Len)
	}

	k := key.NodePrivateFrom([32]byte(b))

	return &k, nil
}

func WriteKey(k key.NodePrivate, keyFile string) error {
	nKey := key.UnveilPrivate(k)
	return os.WriteFile(keyFile, nKey[:], 0o600)
}
