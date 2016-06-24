package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	kubeclient "k8s.io/kubernetes/pkg/client/unversioned"
)

const (
	numPodsQueryKey             = "num_pods"
	numConcurrentPerPodQueryKey = "num_concurrent_per_pod"
	numReqsPerPodQueryKey       = "num_reqs_per_pod"
	httpMethodQueryKey          = "http_method"
	endpointQueryKey            = "endpoint"
)

var (
	allowedHTTPMethods = map[string]struct{}{
		"GET":     struct{}{},
		"POST":    struct{}{},
		"PUT":     struct{}{},
		"DELETE":  struct{}{},
		"OPTIONS": struct{}{},
	}
)

type startReqBody struct {
	NumPods             int    `json:"num_pods"`
	NumConcurrentPerPod int    `json:"num_concurrent_per_pod"`
	NumReqsPerPod       int    `json:"num_reqs_per_pod"`
	HTTPMethod          string `json:"http_method"`
	Endpoint            string `json:"endpoint"`
	Namespace           string `json:"namespace"`
	Image               string `json:"image"`
}

func newStartJobHandler(kcl *kubeclient.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "you must POST to this", http.StatusBadRequest)
			return
		}
		req := new(startReqBody)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, fmt.Sprintf("decoding request body (%s)", err), http.StatusBadRequest)
			return
		}
		boomCmd := boomCommand{total: req.NumReqsPerPod, concurrency: req.NumConcurrentPerPod, endpoint: req.Endpoint}
		boomJob := newBoomJob(req.Image, boomCmd, defaultJobNamespace, defaultJobName, req.NumPods)
		if _, err := kcl.BatchClient.Jobs(req.Namespace).Create(boomJob); err != nil {
			http.Error(w, fmt.Sprintf("error creating job (%s)", err), http.StatusInternalServerError)
			return
		}
	})
}
