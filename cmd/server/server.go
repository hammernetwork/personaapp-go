package server

import (
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var config Config

func Command() *cobra.Command {


	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start a Persona App Server",
		RunE: run,
	}
	cmd.Flags().AddFlagSet(config.Flags())
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
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

	sugar.Info("Environment: ", config.Environment)

	return nil
}
