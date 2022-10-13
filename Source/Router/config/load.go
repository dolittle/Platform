package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"sync"
)

// LoadConfigFor loads the configuration using the given cobra.Command flags,
// and any supplied YAML files through '--config' arguments. It also starts
// file watchers for each of the configuration files to reload config when it
// changes.
func LoadConfigFor(cmd *cobra.Command) (*Config, error) {
	c := &Config{
		k: koanf.New("."),
		l: sync.RWMutex{},
	}

	configFiles, _ := cmd.Flags().GetStringSlice("config")
	for _, configFile := range configFiles {
		provider := file.Provider(configFile)
		if err := c.k.Load(provider, yaml.Parser()); err != nil {
			return nil, err
		}
		_ = provider.Watch(func(_ any, err error) {
			if err != nil {
				log.Warn().Err(err).Msg("Config file watcher failed")
				return
			}

			c.l.Lock()
			defer c.l.Unlock()

			err = c.k.Load(provider, yaml.Parser())
			if err != nil {
				log.Warn().Err(err).Msg("Failed to re-load changed config")
				return
			}

			log.Info().Str("path", configFile).Msg("Re-loaded config from file")
		})
	}

	if err := c.k.Load(posflag.Provider(cmd.Flags(), c.k.Delim(), c.k), nil); err != nil {
		return nil, err
	}

	return c, nil
}
