package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

type flagSet struct {
	fs       *pflag.FlagSet
	required []string
	addFlags func(*pflag.FlagSet)
}

func newFlagSet(addFlags func(*pflag.FlagSet), requiredFields ...string) *flagSet {
	return &flagSet{
		fs:       pflag.NewFlagSet("megaboom", pflag.ContinueOnError),
		required: requiredFields,
		addFlags: addFlags,
	}
}

func parseAndValidate(fs *flagSet, args []string) error {
	fmt.Println("PARSEANDVALIDATE")
	fmt.Println("CALLING")
	fs.addFlags(fs.fs)

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
