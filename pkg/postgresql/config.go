package postgresql

import (
	"errors"

	"github.com/spf13/pflag"
)

type Config struct {
	Host               string
	Port               uint16
	User               string
	Password           string
	Database           string
	MaxOpenConnections int
	MaxIdleConnections int
}

func (c *Config) Flags(name, prefix string) *pflag.FlagSet {
	if prefix != "" {
		prefix += "_"
	}

	f := pflag.NewFlagSet(name, pflag.PanicOnError)
	f.StringVar(&c.Host, prefix+"host", "127.0.0.1", "Host")
	f.Uint16Var(&c.Port, prefix+"port", 5432, "Port")
	f.StringVar(&c.User, prefix+"user", "root", "User")
	f.StringVar(&c.Password, prefix+"password", "", "Password")
	f.StringVar(&c.Database, prefix+"database", "", "Database")
	f.IntVar(&c.MaxOpenConnections, prefix+"max_open_connections", 16, "Max number of the open connections")
	f.IntVar(&c.MaxIdleConnections, prefix+"max_idle_connections", 16, "Max number of the idle connections")

	return f
}

func (c *Config) Validate() error {
	if c.Password == "" || c.Database == "" {
		return errors.New("empty config")
	}

	return nil
}
