package main

import (
	"context"
	"fmt"

	"github.com/arschles/megaboom/pkg/k8s"
	"github.com/google/uuid"
	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type startCommand struct {
	ui            cli.Ui
	fs            *pflag.FlagSet
	endpoint      *string
	reqsPerRunner *uint
	concurrency   *uint
	numRunners    *uint
	ns            *string
	requiredFlags []string
}

func (s startCommand) Help() string {
	return fmt.Sprintf(
		"Start a load testing job with a specified concurrency, total number of requests, and endpoint. Usage:\n%s",
		s.fs.FlagUsages(),
	)
}

func (s startCommand) Synopsis() string {
	return "Start a load testing job with a specified concurrency, total number of requests, and endpoint"
}

func (s startCommand) Run(args []string) int {
	ctx := context.Background()

	if err := parseAndValidate(s.fs, args, s.requiredFlags...); err != nil {
		s.ui.Error(err.Error())
		return 1
	}

	cl, err := k8s.NewClient(false)
	if err != nil {
		s.ui.Error(err.Error())
		return 1
	}

	uid := uuid.New()
	job := k8s.NewJob(
		uid,
		*s.endpoint,
		*s.numRunners,
		*s.reqsPerRunner,
		*s.concurrency,
	)
	jobCreator := cl.BatchV1().Jobs(*s.ns)
	if _, err := jobCreator.Create(ctx, job, metav1.CreateOptions{}); err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	return 0
}

func startCommandFactory(ui cli.Ui) cli.CommandFactory {
	flagSet := pflag.NewFlagSet("megaboom", pflag.ExitOnError)
	endpoint := flagSet.StringP("endpoint", "e", "", "The endpoint to hit")
	reqsPerRunner := flagSet.UintP("reqs-per-runner", "r", 1, "The number of requests to make per runner")
	concurrency := flagSet.UintP("concurrency", "c", 1, "The number of concurrent requests to make")
	numRunners := flagSet.UintP("num-runners", "t", 1, "The total number of runners (pods) to run")
	ns := flagSet.StringP("namespace", "n", "default", "The namespace to run the load test in")

	return func() (cli.Command, error) {
		return startCommand{
			ui:            ui,
			fs:            flagSet,
			endpoint:      endpoint,
			reqsPerRunner: reqsPerRunner,
			concurrency:   concurrency,
			numRunners:    numRunners,
			ns:            ns,
			requiredFlags: []string{"endpoint"},
		}, nil
	}
}
