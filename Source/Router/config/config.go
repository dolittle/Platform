package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"sync"
)

type Config struct {
	k *koanf.Koanf
	l sync.RWMutex
	c chan struct{}
}

func (c *Config) Changed() <-chan struct{} {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.c
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

func (c *Config) Unmarshal(path string, o any) error {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.k.Unmarshal(path, o)
}

func (c *Config) MarshalYaml() ([]byte, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.k.Marshal(yaml.Parser())
}
