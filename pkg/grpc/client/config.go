package client

import (
	"errors"
	"time"

	"personaapp/pkg/flag/mapping"

	"github.com/spf13/pflag"
)

type Config struct {
	Servers        string
	DialTimeout    time.Duration
	RequestTimeout time.Duration
}

func (c *Config) Flags(prefix string) *pflag.FlagSet {
	name := prefix + "Client"
	f := pflag.NewFlagSet(name, pflag.PanicOnError)

	f.StringVar(&c.Servers, "servers", "", "ip:port address of the service")
	f.DurationVar(&c.DialTimeout, "dial_timeout", 30*time.Second, "service client dial timeout")
	f.DurationVar(&c.RequestTimeout, "request_timeout", 60*time.Second, "service client request timeout")

	return mapping.WithPrefix(f, name, pflag.PanicOnError, prefix)
}

func (c *Config) Validate() error {
	if c.Servers == "" {
		return errors.New("servers is not configured")
	}

	return nil
}
