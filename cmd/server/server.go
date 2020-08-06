package server

import (
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"personaapp/internal/server"
	"personaapp/pkg/closeable"

	"github.com/improbable-eng/grpc-web/go/grpcweb"

	authController "personaapp/internal/controllers/auth/controller"
	authStorage "personaapp/internal/controllers/auth/storage"
	companyController "personaapp/internal/controllers/company/controller"
	companyStorage "personaapp/internal/controllers/company/storage"
	vacancyController "personaapp/internal/controllers/vacancy/controller"
	vacancyStorage "personaapp/internal/controllers/vacancy/storage"
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

		srv := server.New(
			newAuthController(pg, &cfg.AuthController),
			newCompanyController(pg),
			newVacancyController(pg),
		)

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
		startHTTPServer(grpcServer)

		pkgcmd.Await()
		grpcServer.GracefulStop()
		if err := g.Wait(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}

func startHTTPServer(grpcServer *grpc.Server) {
	grpcWebServer := grpcweb.WrapServer(grpcServer)
	_ = &http.Server{
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, "+
					"Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
	log.Printf("starting HTTP Server")
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
