package server

import (
	"log"
	"net"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"personaapp/internal/server"
	"personaapp/internal/server/controller"
	"personaapp/internal/server/storage"
	"personaapp/pkg/closeable"
	pkgcmd "personaapp/pkg/cmd"
	"personaapp/pkg/flag"
	"personaapp/pkg/grpcapi/personaappapi"
	"personaapp/pkg/postgresql"
)

func Command() *cobra.Command {
	var config Config
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a Persona App Server",
		RunE:  run(&config),
	}
	cmd.Flags().AddFlagSet(config.Flags())
	return cmd
}

func run(cfg *Config) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		logger, _ := zap.NewProduction()
		defer func() {
			err := logger.Sync() // flushes buffer, if any
			if err != nil {
				log.Println(err) // todo think about errors mapper/parser service
			}
		}()
		sugar := logger.Sugar()

		sugar.Info("starting server")
		defer sugar.Info("stopping server")

		if err := flag.BindEnv(cmd); err != nil {
			return errors.WithStack(err)
		}

		pg, err := postgresql.New(&cfg.Postgres)
		if err != nil {
			return errors.WithStack(err)
		}
		closeable.CloseWithErrorLogging(sugar, pg)

		s := storage.New(pg)
		c := controller.New(s)
		srv := server.New(c)

		ln, err := net.Listen("tcp", cfg.Server.Address)
		if err != nil {
			return errors.WithStack(err)
		}

		grpcServer := grpc.NewServer()
		personaappapi.RegisterPersonaAppServer(grpcServer, srv)
		reflection.Register(grpcServer)

		g := &errgroup.Group{}
		g.Go(func() error {
			if err := grpcServer.Serve(ln); err != nil {
				return errors.WithStack(err)
			}
			return nil
		})

		pkgcmd.Await()
		grpcServer.GracefulStop()
		if err := g.Wait(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}
