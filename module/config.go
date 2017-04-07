package module

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type Config struct {
	DSN   string   `toml:"dsn"`
	Token string   `toml:"token"`
	Root  []string `toml:"root"`
}

func (c *Config) Parse(path string) (err error) {
	_, err = toml.DecodeFile(path, c)
	if err != nil {
		err = errors.Wrap(err, "Config.Read error")
	}
	return
}
