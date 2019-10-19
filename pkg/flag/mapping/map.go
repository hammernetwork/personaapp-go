package mapping

import "github.com/spf13/pflag"

func WithPrefix(f *pflag.FlagSet, name string, errorHandling pflag.ErrorHandling, prefix string) *pflag.FlagSet {
	return all(f, name, errorHandling, func(flag *pflag.Flag) {
		if prefix != "" {
			flag.Name = prefix + "." + flag.Name
		}
	})
}

func all(f *pflag.FlagSet, name string, errorHandling pflag.ErrorHandling, fn func(*pflag.Flag)) *pflag.FlagSet {
	fNew := pflag.NewFlagSet(name, errorHandling)

	f.VisitAll(func(flag *pflag.Flag) {
		fn(flag)
		fNew.AddFlag(flag)
	})

	return fNew
}
