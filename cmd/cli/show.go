package main

import (
	"context"
	"fmt"

	"github.com/arschles/megaboom/pkg/k8s"
	"github.com/mitchellh/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type showCommand struct {
	ui   cli.Ui
	fs   *flagSet
	name string
	ns   string
}

func (s showCommand) Help() string {
	return fmt.Sprintf(
		"Show details on a load testing job. Usage:\n%s",
		s.fs.fs.FlagUsages(),
	)
}

func (s showCommand) Synopsis() string {
	return "Show details on a load testing job"
}

func (s showCommand) Run(args []string) int {
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

	// get details about the job and also the config map with reports
	// (if the job is done)
	jobGetter := cl.BatchV1().Jobs(s.ns)
	job, err := jobGetter.Get(ctx, s.name, metav1.GetOptions{})
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	s.ui.Info(fmt.Sprintf(
		"Job succeeded? %t",
		job.Status.Succeeded == 1,
	))
	s.ui.Info(fmt.Sprintf(
		"%d workers succeeded / %d total workers",
		job.Status.Succeeded,
		*job.Spec.Completions,
	))
	uid, err := k8s.GetUIDFromJob(job)
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	endpoint, err := k8s.GetEndpointFromJob(job)
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	s.ui.Info(fmt.Sprintf("UID: %s", uid.String()))
	s.ui.Info(fmt.Sprintf("Endpoint: %s", endpoint))
	return 0
}

func showCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		fs := newFlagSet("name")
		cmd := showCommand{ui: ui, fs: fs}

		fs.fs.StringVarP(&cmd.ns, "namespace", "n", "default", "The namespace to run the load test in")
		fs.fs.StringVarP(&cmd.name, "name", "n", "", "The name of the load test job")
		// TODO: watch
		return cmd, nil
	}
}
