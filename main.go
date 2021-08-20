package main

import (
	"fmt"
	"os"

	"github.com/arschles/megaboom/pkg/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Error creating logger: ", err)
		os.Exit(1)
	}
	lggr := zapr.NewLogger(zapLogger)

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
	r.GET("/healthz", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.String(200, "OK")
	})
	r.POST("/job", handlers.StartJob(lggr, kcl))
	r.DELETE("/job/:id", handlers.DeleteJob(lggr, kcl))

	lggr.Info("starting megaboom server", "port", "8080")
	lggr.Error(r.Run(":8080"), "server failed")
}
