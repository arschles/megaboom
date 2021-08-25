package main

import (
	"context"
	"fmt"

	"github.com/arschles/megaboom/pkg/k8s"
	"github.com/mitchellh/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type listCommand struct {
	ui cli.Ui
	fs *flagSet
	ns string
}

func (s listCommand) Help() string {
	return fmt.Sprintf(
		"Start a load testing job with a specified concurrency, total number of requests, and endpoint. Usage:\n%s",
		s.fs.fs.FlagUsages(),
	)
}

func (s listCommand) Synopsis() string {
	return "Start a load testing job with a specified concurrency, total number of requests, and endpoint"
}

func (s listCommand) Run(args []string) int {
	ctx := context.Background()

	if err := parseAndValidate(s.fs, args); err != nil {
		s.ui.Error(err.Error())
		return 1
	}

	cl, err := k8s.NewClient(false)
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}

	jobLister := cl.BatchV1().Jobs(s.ns)
	jobs, err := jobLister.List(ctx, metav1.ListOptions{})
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	s.ui.Info(fmt.Sprintf(
		"jobs in namespace %s:\n", s.ns,
	))
	for _, job := range jobs.Items {
		s.ui.Info(job.Name)
	}
	return 0
}

func listCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		fs := newFlagSet()
		cmd := listCommand{ui: ui, fs: fs}
		fs.fs.StringVarP(&cmd.ns, "namespace", "n", "default", "namespace to list jobs in")
		return cmd, nil
	}
}
