package cmd

import (
	"personaapp/cmd/server"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

func Run() error {
	rootCMD := &cobra.Command{}
	rootCMD.AddCommand(server.Command())

	return errors.WithStack(rootCMD.Execute())
}
