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
	fs            *flagSet
	endpoint      string
	reqsPerRunner uint
	concurrency   uint
	numRunners    uint
	ns            string
	headers       []string
}

func (s startCommand) Help() string {
	s.fs.addFlags(s.fs.fs)
	return fmt.Sprintf(
		"Start a load testing job with a specified concurrency, total number of requests, and endpoint. Usage:\n%s",
		s.fs.fs.FlagUsages(),
	)
}

func (s startCommand) Synopsis() string {
	return "Start a load testing job with a specified concurrency, total number of requests, and endpoint"
}

func (s startCommand) Run(args []string) int {
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

	uid := uuid.New()
	job := k8s.NewJob(
		uid,
		s.endpoint,
		s.numRunners,
		s.reqsPerRunner,
		s.concurrency,
		s.headers,
	)
	s.ui.Info(fmt.Sprintf(
		"Creating job in namespace %s with uid %s",
		s.ns,
		uid.String(),
	))
	jobCreator := cl.BatchV1().Jobs(s.ns)
	if _, err := jobCreator.Create(ctx, job, metav1.CreateOptions{}); err != nil {
		s.ui.Error(err.Error())
		return 1
	}
	return 0
}

func startCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		cmd := startCommand{ui: ui}
		fs := newFlagSet(func(fs *pflag.FlagSet) {
			fmt.Println("FUYNC")
			fs.StringVarP(&cmd.endpoint, "endpoint", "e", "", "The endpoint to hit")
			fs.UintVarP(&cmd.reqsPerRunner, "reqs-per-runner", "r", 1, "The number of requests to make per runner")
			fs.UintVarP(&cmd.concurrency, "concurrency", "c", 1, "The number of concurrent requests to make per runner")
			fs.UintVarP(&cmd.numRunners, "num-runners", "t", 1, "The total number of runners (pods) to run")
			fs.StringVarP(&cmd.ns, "namespace", "n", "default", "The namespace to run the load test in")
			fs.StringSliceVarP(&cmd.headers, "headers", "h", []string{}, "The headers to send with each request")
		}, "endpoint")
		cmd.fs = fs
		return cmd, nil
	}
}
