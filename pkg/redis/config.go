package redis

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/pflag"
)

type Config redis.Options
type ClusterOptions redis.ClusterOptions

func (c *Config) Flags(prefix string) *pflag.FlagSet { //nolint:funlen
	f := pflag.NewFlagSet("RedisConfig", pflag.PanicOnError)

	if prefix != "" {
		prefix += "."
	}

	f.StringVar(
		&c.Network,
		prefix+"network",
		"tcp",
		"network type, either tcp or unix",
	)
	f.StringVar(
		&c.Addr,
		prefix+"addr",
		"",
		"host:port address",
	)
	f.StringVar(
		&c.Password,
		prefix+"password",
		"",
		"optional password. must match the password specified in the requirepass server configuration option.",
	)
	f.IntVar(
		&c.DB,
		prefix+"db",
		0,
		"db id",
	)
	f.IntVar(
		&c.MaxRetries,
		prefix+"max_retries",
		0,
		"maximum number of retries before giving up. default is to not retry failed commands",
	)
	f.DurationVar(
		&c.MinRetryBackoff,
		prefix+"min_retry_backoff",
		8*time.Millisecond,
		"minimum backoff between each retry. -1 disables backoff.",
	)
	f.DurationVar(
		&c.MaxRetryBackoff,
		prefix+"max_retry_backoff",
		512*time.Millisecond,
		"maximum backoff between each retry. -1 disables backoff.",
	)
	f.DurationVar(
		&c.DialTimeout,
		prefix+"dial_timeout",
		5*time.Second,
		"dial timeout for establishing new connections",
	)
	f.DurationVar(
		&c.ReadTimeout,
		prefix+"read_timeout",
		3*time.Second,
		"timeout for socket reads. If reached, commands will fail with a timeout instead of blocking",
	)
	f.DurationVar(
		&c.WriteTimeout,
		prefix+"write_timeout",
		3*time.Second,
		"timeout for socket writes. If reached, commands will fail with a timeout instead of blocking",
	)
	f.IntVar(
		&c.PoolSize,
		prefix+"pool_size",
		10,
		"maximum number of socket connections",
	)
	f.DurationVar(
		&c.PoolTimeout,
		prefix+"pool_timeout",
		4*time.Second,
		"amount of time client waits for connection if all connections are busy before returning an error",
	)
	f.DurationVar(
		&c.IdleTimeout,
		prefix+"idle_timeout",
		5*time.Minute,
		"amount of time after which client closes idle connections. Should be less than server's timeout",
	)
	f.DurationVar(
		&c.IdleCheckFrequency,
		prefix+"idle_check_frequency",
		1*time.Minute,
		"frequency of idle checks. when minus value is set, then idle check is disabled.",
	)

	return f
}

func (c *ClusterOptions) Flags(prefix string) *pflag.FlagSet { //nolint:funlen
	f := pflag.NewFlagSet("RedisClusterConfig", pflag.PanicOnError)

	if prefix != "" {
		prefix += "."
	}

	f.StringArrayVar(
		&c.Addrs,
		prefix+"addrs",
		[]string{},
		"a seed list of host:port addresses of cluster nodes.",
	)
	f.StringVar(
		&c.Password,
		prefix+"password",
		"",
		"optional password. must match the password specified in the requirepass server configuration option.",
	)
	f.IntVar(
		&c.MaxRetries,
		prefix+"max_retries",
		0,
		"maximum number of retries before giving up. default is to not retry failed commands.",
	)
	f.DurationVar(
		&c.DialTimeout,
		prefix+"dial_timeout",
		5*time.Second,
		"dial timeout for establishing new connections.",
	)
	f.DurationVar(
		&c.ReadTimeout,
		prefix+"read_timeout",
		3*time.Second,
		"timeout for socket reads. if reached, commands will fail with a timeout instead of blocking.",
	)
	f.DurationVar(
		&c.WriteTimeout,
		prefix+"write_timeout",
		3*time.Second,
		"timeout for socket writes. if reached, commands will fail with a timeout instead of blocking.",
	)
	f.IntVar(
		&c.PoolSize,
		prefix+"pool_size",
		10, "maximum number of socket connections",
	)
	f.DurationVar(
		&c.PoolTimeout,
		prefix+"pool_timeout",
		4*time.Second,
		"amount of time client waits for connection if all connections are busy before returning an error.",
	)
	f.DurationVar(
		&c.IdleTimeout,
		prefix+"idle_timeout",
		5*time.Minute,
		"amount of time after which client closes idle connections. should be less than server's timeout.",
	)
	f.DurationVar(
		&c.IdleCheckFrequency,
		prefix+"idle_check_frequency",
		1*time.Minute, "frequency of idle checks. when minus value is set, then idle check is disabled.",
	)
	f.BoolVar(
		&c.ReadOnly,
		prefix+"readonly",
		false,
		"enables read only queries on slave nodes.",
	)
	f.BoolVar(
		&c.RouteByLatency,
		prefix+"route_by_latency",
		false,
		"enables routing read-only queries to the closest master or slave node.",
	)

	return f
}
