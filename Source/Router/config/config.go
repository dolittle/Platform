package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"sync"
)

type Config struct {
	k *koanf.Koanf
	l sync.RWMutex
}

func (c *Config) String(path string) string {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.k.String(path)
}

func (c *Config) Int(path string) int {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.k.Int(path)
}

func (c *Config) MarshalYaml() ([]byte, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.k.Marshal(yaml.Parser())
}
