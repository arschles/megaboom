package main

import (
	"context"
	"fmt"

	"github.com/arschles/megaboom/pkg/k8s"
	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deleteCommand struct {
	ui cli.Ui
}

func (s deleteCommand) Help() string {
	return "Delete a load testing job and its reports"
}

func (s deleteCommand) Synopsis() string {
	return "Delete an existing load testing job and its reports"
}

func (s deleteCommand) Run(args []string) int {
	ctx := context.Background()
	flagSet := pflag.NewFlagSet("megaboom", pflag.ExitOnError)
	ns := flagSet.StringP("namespace", "n", "default", "The namespace to run the load test in")
	name := flagSet.StringP("name", "", "", "The name of the load test job to delete")

	if err := flagSet.Parse(args); err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	cl, err := k8s.NewClient(false)
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}

	jobDeleter := cl.BatchV1().Jobs(*ns)
	if err := jobDeleter.Delete(ctx, *name, metav1.DeleteOptions{}); err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	s.ui.Info(fmt.Sprintf(
		"deleted load testing job %s/%s", *ns, *name),
	)
	return 0
}

func deleteCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return deleteCommand{ui: ui}, nil
	}
}
