package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func parseAndValidate(fs *pflag.FlagSet, args []string, requiredFlags ...string) error {
	if err := fs.Parse(args); err != nil {
		return err
	}
	return requireFlags(fs, requiredFlags...)
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
