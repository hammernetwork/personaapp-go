package redis

import "github.com/go-redis/redis"

type DefaultStorage struct {
	Client *redis.Client
}

type ClusterStorage struct {
	Client *redis.ClusterClient
}

func NewStorage(conf *Config) (*DefaultStorage, error) {
	s := &DefaultStorage{
		Client: redis.NewClient((*redis.Options)(conf)),
	}

	return s, nil
}

func NewClusterStorage(conf *ClusterOptions) (*ClusterStorage, error) {
	s := &ClusterStorage{
		Client: redis.NewClusterClient((*redis.ClusterOptions)(conf)),
	}

	return s, nil
}
