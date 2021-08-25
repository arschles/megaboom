package main

import (
	"context"
	"fmt"
	"os"

	"github.com/arschles/megaboom/pkg/k8s"
	"github.com/arschles/megaboom/pkg/log"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	lggr, err := log.NewZapr()
	if err != nil {
		fmt.Printf(
			"couldn't create logger (%s)",
			err,
		)
		os.Exit(1)
	}
	cfg := new(config)
	if err := envconfig.Process("MEGABOOM_WORKER", cfg); err != nil {
		lggr.Error(
			err,
			"couldn't process config",
		)
		os.Exit(1)
	}
	ctx, done := context.WithTimeout(
		context.Background(),
		cfg.MaxRunTime,
	)
	defer done()
	runnerID := uuid.New()
	kcl, err := k8s.NewK8sClient()
	if err != nil {
		lggr.Error(
			err,
			"couldn't create Kubernetes client",
		)
		os.Exit(1)
	}
	lggr.Info(
		"starting run",
		[]interface{}{
			"runnerID",
			runnerID.String(),
			cfg.loggerVals(),
		}...,
	)
	result, err := run(ctx, lggr, cfg.TotalRequests, cfg.Concurrency, cfg.Endpoint)
	if err != nil {
		lggr.Error(
			err,
			"error running requests",
		)
		os.Exit(1)
	}
	cm, err := configMapFromResult("megaboom-results", runnerID, result)
	if err != nil {
		lggr.Error(
			err,
			"couldn't convert result to a config map",
		)
		os.Exit(1)
	}
	if _, err := kcl.CoreV1().ConfigMaps(cfg.Namespace).Create(
		ctx,
		&cm,
		metav1.CreateOptions{},
	); err != nil {
		lggr.Error(
			err,
			"couldn't save config map",
		)
		os.Exit(1)
	}

}
