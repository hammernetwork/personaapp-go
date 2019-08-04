package grpc

import (
	"time"

	"personaapp/pkg/flag/mapping"

	"github.com/spf13/pflag"
)

type Config struct {
	Address        string
	RequestTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func (c *Config) Flags(prefix string) *pflag.FlagSet {
	name := prefix + "grpc"
	f := pflag.NewFlagSet(name, pflag.PanicOnError)

	f.StringVar(&c.Address, "address", "127.0.0.1:8000", "Address in ip:port format")
	f.DurationVar(&c.RequestTimeout, "request_timeout", 60*time.Second, "Request timeout")
	f.DurationVar(&c.ReadTimeout, "read_timeout", 60*time.Second, "Read timeout")
	f.DurationVar(&c.WriteTimeout, "write_timeout", 60*time.Second, "Write timeout")

	return mapping.WithPrefix(f, name, pflag.PanicOnError, prefix)
}
