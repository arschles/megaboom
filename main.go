package main

import (
	"flag"
	"log"

	// "k8s.io/kubernetes/pkg/api"
	kubeclient "k8s.io/kubernetes/pkg/client/unversioned"
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

func main() {
	numPods := flag.Int("p", 10, "Number of pods to run. Each pod will run the same command, based on the other flags given.")
	numReqs := flag.Int("n", 10, "Number of requests to run.")
	numConcurrent := flag.Int("c", 1, "Number of requests to run concurrently. Total number of requests cannot be smaller than the concurrency level.")
	// not yet implemented
	// -q  Rate limit, in seconds (QPS).
	// not yet implemented
	// -o  Output type. If none provided, a summary is printed.
	//     "csv" is the only supported alternative. Dumps the response
	//     metrics in comma-seperated values format.
	httpMethod := flag.String("m", "GET", "HTTP method, one of GET, POST, PUT, DELETE, HEAD, OPTIONS.")
	// not yet implemented
	// -H  Custom HTTP header. You can specify as many as needed by repeating the flag.
	// for example, -H "Accept: text/html" -H "Content-Type: application/xml" .
	timeoutMS := flag.Int("t", 100, "Timeout in ms.")
	// not yet implemented
	// -A  HTTP Accept header.
	// not yet implemented
	// -d  HTTP request body.
	// not yet implemented
	// -T  Content-type, defaults to "text/html".
	// not yet implemented
	// -a  Basic authentication, username:password.
	// not yet implemented
	// -x  HTTP Proxy address as host:port.

	// not yet implemented
	// -disable-compression  Disable compression.
	// not yet implemented
	// -disable-keepalive    Disable keep-alive, prevents re-use of TCP
	// connections between different HTTP requests.
	// not yet implemented
	// -cpus                 Number of used cpu cores.
	// (default for current machine is 1 cores)
	flag.Parse()
	if *numConcurrent > *numReqs {
		log.Fatalf("%d concurrent requests is greater than %d total requests", *numConcurrent, *numReqs)
	}
	if _, ok := allowedHTTPMethods[*httpMethod]; !ok {
		log.Fatalf("%s is an unallowed HTTP method", *httpMethod)
	}
	if *numPods < 1 {
		log.Fatalf("There must be at least 1 pod to run")
	}

	kcl, err := kubeclient.NewInCluster()
	if err != nil {
		log.Fatalf("Error creating new Kubernetes client (%s)", err)
	}
	boomJob := newBoomJob(defaultBoomImage)
	kcl.Jobs(namespace).Create(boomJob)

}
