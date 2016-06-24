package main

import (
	"log"
	"net/http"

	kubeclient "k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {
	kcl, err := kubeclient.NewInCluster()
	if err != nil {
		log.Fatalf("Error creating new Kubernetes client (%s)", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/job/create", newStartJobHandler(kcl))
	log.Printf("serving on 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("error serving (%s)", err)
	}
}
