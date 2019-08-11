package server

import (
	"github.com/spf13/pflag"

	"personaapp/pkg/grpc"
	"personaapp/pkg/postgresql"
)

type Config struct {
	Postgres    postgresql.Config
	Server      grpc.Config
	Environment string
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("ServerConfig", pflag.PanicOnError)

	f.AddFlagSet(c.Postgres.Flags("PostgresConfig", "postgres"))
	f.AddFlagSet(c.Server.Flags("ServerConfig", "server"))
	f.StringVar(&c.Environment, "environment", "dev", "Test environment variable")

	return f
}
