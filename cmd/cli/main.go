package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
)

type config struct {
	Namespace string `envconfig:"NAMESPACE" required:"true"`
}

func main() {
	ui := cli.BasicUi{
		Writer:      os.Stdout,
		Reader:      os.Stdin,
		ErrorWriter: os.Stderr,
	}
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"start": startCommandFactory(&ui),
		"list":  listCommandFactory(&ui),
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	os.Exit(exitStatus)

	// zapLogger, err := zap.NewDevelopment()
	// if err != nil {
	// 	fmt.Println("Error creating logger: ", err)
	// 	os.Exit(1)
	// }
	// lggr := zapr.NewLogger(zapLogger)

	// cfg := new(config)
	// if err := envconfig.Process("MEGABOOM", cfg); err != nil {
	// 	lggr.Error(err, "couldn't process config, exiting")
	// 	os.Exit(1)
	// }

	// kcl, err := k8s.NewK8sClient()
	// if err != nil {
	// 	lggr.Error(err, "failed to create new kube client")
	// 	os.Exit(1)
	// }

	// r := gin.Default()
	// r.GET("/livez", func(c *gin.Context) {
	// 	c.String(200, "OK")
	// })
	// r.GET("/readyz", func(c *gin.Context) {
	// 	c.String(200, "OK")
	// })
	// r.POST("/job", handlers.StartJob(lggr, kcl.BatchV1().Jobs(cfg.Namespace)))
	// r.DELETE("/job/:id", handlers.DeleteJob(lggr, kcl.BatchV1().Jobs(cfg.Namespace)))

	// lggr.Info("starting megaboom server", "port", "8080")
	// lggr.Error(r.Run(":8080"), "server failed")
}
