package client

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func New(
	ctx context.Context,
	cfg *Config,
	bind func(*grpc.ClientConn) interface{},
	customOpts ...Option,
) (interface{}, *grpc.ClientConn, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, errors.WithStack(err)
	}

	// block execution until the connection is established
	opt := option{
		block: true,
	}

	for _, o := range customOpts {
		o(&opt)
	}

	var interceptors = []grpc.UnaryClientInterceptor{
		grpc_prometheus.UnaryClientInterceptor,
		// prometheus.UnaryClientInterceptor, TODO: client metrics
	}

	if opt.unaryInt != nil {
		interceptors = append(interceptors, *opt.unaryInt)
	}

	dialOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)),
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	}
	if opt.block {
		dialOpts = append(dialOpts, grpc.WithBlock())
	}

	if cfg.DialTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.DialTimeout)

		defer cancel()
	}

	conn, err := grpc.DialContext(ctx, cfg.Servers, dialOpts...)

	if err != nil {
		return nil, nil, err
	}

	return bind(conn), conn, nil
}

type option struct {
	unaryInt *grpc.UnaryClientInterceptor
	block    bool
}

type Option func(*option)

func WithCustomUnaryClientInterceptor(interceptor *grpc.UnaryClientInterceptor) Option {
	return func(option *option) {
		option.unaryInt = interceptor
	}
}

// WithAsync instructs dial to connect in a non-blocking mod
func WithAsync() Option {
	return func(option *option) {
		option.block = false
	}
}
