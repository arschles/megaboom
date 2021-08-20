package main

import (
	"fmt"
	"os"

	"github.com/arschles/megaboom/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/zapr"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type config struct {
	Namespace string `envconfig:"NAMESPACE" required:"true"`
}

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Error creating logger: ", err)
		os.Exit(1)
	}
	lggr := zapr.NewLogger(zapLogger)

	cfg := new(config)
	if err := envconfig.Process("MEGABOOM", cfg); err != nil {
		lggr.Error(err, "couldn't process config, exiting")
		os.Exit(1)
	}

	restCfg, err := rest.InClusterConfig()
	if err != nil {
		lggr.Error(err, "failed to get in-cluster Kubernetes config")
		os.Exit(1)
	}

	kcl, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		lggr.Error(err, "failed to create new kube client")
		os.Exit(1)
	}

	r := gin.Default()
	r.GET("/livez", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.POST("/job", handlers.StartJob(lggr, kcl.BatchV1().Jobs(cfg.Namespace)))
	r.DELETE("/job/:id", handlers.DeleteJob(lggr, kcl.BatchV1().Jobs(cfg.Namespace)))

	lggr.Info("starting megaboom server", "port", "8080")
	lggr.Error(r.Run(":8080"), "server failed")
}
