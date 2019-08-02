package server

import "github.com/spf13/pflag"

type Config struct {
	Environment string
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("ServerConfig", pflag.PanicOnError)

	f.StringVar(&c.Environment, "environment", "dev", "Test environment variable")

	return f
}

