package server

import (
	"log"
	"net"
	"personaapp/pkg/closeable"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"personaapp/internal/server"
	authController "personaapp/internal/server/auth/controller"
	authStorage "personaapp/internal/server/auth/storage"
	companyController "personaapp/internal/server/company/controller"
	companyStorage "personaapp/internal/server/company/storage"
	pkgcmd "personaapp/pkg/cmd"
	"personaapp/pkg/flag"
	apiauth "personaapp/pkg/grpcapi/auth"
	apicompany "personaapp/pkg/grpcapi/company"
	apivacancy "personaapp/pkg/grpcapi/vacancy"
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
				log.Println(err) //nolint todo think about errors mapper/parser service
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
		// nolint TODO: not sure if there should be defer, but I guess so
		defer closeable.CloseWithErrorLogging(sugar, pg)

		as := authStorage.New(pg)
		ac := authController.New(&cfg.AuthController, as)

		cs := companyStorage.New(pg)
		cc := companyController.New(cs)

		srv := server.New(ac, cc)

		ln, err := net.Listen("tcp", cfg.Server.Address)
		if err != nil {
			return errors.WithStack(err)
		}

		grpcServer := grpc.NewServer()
		apiauth.RegisterPersonaAppAuthServer(grpcServer, srv)
		apicompany.RegisterPersonaAppCompanyServer(grpcServer, srv)
		apivacancy.RegisterPersonaAppVacancyServer(grpcServer, srv)
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
