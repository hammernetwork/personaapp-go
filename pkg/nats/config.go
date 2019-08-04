package nats

import (
	"github.com/nats-io/go-nats"
	"github.com/spf13/pflag"
)

type Config struct {
	Addr string
}

func (c *Config) Flags(prefix string) *pflag.FlagSet {
	f := pflag.NewFlagSet("NatsConfig", pflag.PanicOnError)

	if prefix != "" {
		prefix += "."
	}

	f.StringVar(&c.Addr, prefix+"addr", nats.DefaultURL, "addr of the server to connect")

	return f
}
