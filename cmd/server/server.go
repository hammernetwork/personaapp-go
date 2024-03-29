package server

import (
	"log"
	"net"
	"personaapp/internal/server"
	"personaapp/pkg/closeable"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authController "personaapp/internal/controllers/auth/controller"
	authStorage "personaapp/internal/controllers/auth/storage"
	cityController "personaapp/internal/controllers/city/controller"
	cityStorage "personaapp/internal/controllers/city/storage"
	companyController "personaapp/internal/controllers/company/controller"
	companyStorage "personaapp/internal/controllers/company/storage"
	cvController "personaapp/internal/controllers/cv/controller"
	cvStorage "personaapp/internal/controllers/cv/storage"
	vacancyController "personaapp/internal/controllers/vacancy/controller"
	vacancyStorage "personaapp/internal/controllers/vacancy/storage"
	pkgcmd "personaapp/pkg/cmd"
	"personaapp/pkg/flag"
	apiauth "personaapp/pkg/grpcapi/auth"
	apicity "personaapp/pkg/grpcapi/city"
	apicompany "personaapp/pkg/grpcapi/company"
	apicv "personaapp/pkg/grpcapi/cv"
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

		srv := createControllers(pg, cfg)

		ln, err := net.Listen("tcp", cfg.Server.Address)
		if err != nil {
			return errors.WithStack(err)
		}

		grpcServer := grpc.NewServer()
		registerServer(grpcServer, srv)

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

func registerServer(grpcServer *grpc.Server, srv *server.Server) {
	apiauth.RegisterPersonaAppAuthServer(grpcServer, srv)
	apicompany.RegisterPersonaAppCompanyServer(grpcServer, srv)
	apivacancy.RegisterPersonaAppVacancyServer(grpcServer, srv)
	apicity.RegisterPersonaAppCityServer(grpcServer, srv)
	apicv.RegisterPersonaAppCVServer(grpcServer, srv)
	reflection.Register(grpcServer)
}

func createControllers(pg *postgresql.Storage, cfg *Config) *server.Server {
	return server.New(
		newAuthController(pg, &cfg.AuthController),
		newCompanyController(pg),
		newVacancyController(pg),
		newCityController(pg),
		newCVController(pg),
	)
}

func newAuthController(pg *postgresql.Storage, cfg *authController.Config) *authController.Controller {
	return authController.New(cfg, authStorage.New(pg))
}

func newCompanyController(pg *postgresql.Storage) *companyController.Controller {
	return companyController.New(companyStorage.New(pg))
}

func newVacancyController(pg *postgresql.Storage) *vacancyController.Controller {
	return vacancyController.New(vacancyStorage.New(pg))
}

func newCityController(pg *postgresql.Storage) *cityController.Controller {
	return cityController.New(cityStorage.New(pg))
}

func newCVController(pg *postgresql.Storage) *cvController.Controller {
	return cvController.New(cvStorage.New(pg))
}
