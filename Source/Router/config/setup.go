package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// SetupFor loads the koanf.Koanf configuration and creates a zerolog.Logger for the supplied
// cobra.Command
func SetupFor(cmd *cobra.Command) (*Config, zerolog.Logger, error) {
	config, err := LoadConfigFor(cmd)
	if err != nil {
		return nil, zerolog.Nop(), err
	}

	logger := CreateLoggerUsing(cmd, config)

	return config, logger, nil
}
