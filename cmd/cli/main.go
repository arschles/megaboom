package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	ui := cli.BasicUi{
		Writer:      os.Stdout,
		Reader:      os.Stdin,
		ErrorWriter: os.Stderr,
	}
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"start":  startCommandFactory(&ui),
		"list":   listCommandFactory(&ui),
		"delete": deleteCommandFactory(&ui),
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	os.Exit(exitStatus)
}
