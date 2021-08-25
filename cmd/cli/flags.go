package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

type flagSet struct {
	fs       *pflag.FlagSet
	required []string
}

func newFlagSet(requiredFields ...string) *flagSet {
	return &flagSet{
		fs:       pflag.NewFlagSet("megaboom", pflag.ContinueOnError),
		required: requiredFields,
	}
}

func parseAndValidate(fs *flagSet, args []string) error {
	if err := fs.fs.Parse(args); err != nil {
		return err
	}
	return requireFlags(fs.fs, fs.required...)
}

func requireFlags(fs *pflag.FlagSet, flagNames ...string) error {
	for _, flagName := range flagNames {
		flag := fs.Lookup(flagName)
		if !flag.Changed {
			return fmt.Errorf("flag '%s' not set", flagName)
		}
	}
	return nil
}
