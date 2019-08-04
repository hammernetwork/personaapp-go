package flag

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func BindEnv(cmd *cobra.Command) {
	names := map[string]struct{}{}
	cmd.Flags().Visit(func(f *pflag.Flag) {
		names[f.Name] = struct{}{}
	})

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		replacer := strings.NewReplacer("-", "_", ".", "_")
		name := replacer.Replace(strings.ToUpper(f.Name))

		val := os.Getenv(name)
		if val == "" {
			return
		}

		if _, ok := names[f.Name]; ok {
			return
		}

		t := f.Value.Type()
		if t == "stringArray" || t == "stringSlice" {
			vals := strings.Split(val, " ")
			for _, v := range vals {
				cmd.Flags().Set(f.Name, v)
			}

			return
		}

		cmd.Flags().Set(f.Name, val)
	})
}
