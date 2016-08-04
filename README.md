# megaboom

Megaboom is a Distributed HTTP Load Generator, based on Kubernetes. It is an HTTP server meant to be run inside your Kubernetes cluster, behind a [Service](http://kubernetes.io/docs/user-guide/services/), and it runs [Job](https://github.com/rakyll/boom)s that invocate the [boom](https://github.com/rakyll/boom) CLI to generate load.

# Alpha Status
This project is still in alpha, and is missing a few major features:

- **Reporting the success/failure rate of ongoing jobs** - currently, it's expected that you introspect each [Pod](http://kubernetes.io/docs/user-guide/pods/) to get individual `boom` command output, or use your own monitoring/metrics system to see status of your application under load
- **Cleaning up after itself** - currently it creates a `megaboom` job in the `default` namespace. You must delete that job before you make the next HTTP request to megaboom. Do the deletion with `kubectl delete job megaboom --namespace=${JOB_NAMESPACE}` (see the "Usage" section below for details on the `JOB_NAMESPACE` environment variable)
- **CI infrastructure to run tests and automatically build images** - currently, there are very few tests, and the current Docker image has been manually pushed to `quay.io/arschles/megaboom:devel`

# Usage

This server must run inside a Kubernetes pod. It's recommended that you run this server in a pod managed by a replication controller, and then run a service in front of the replication controller. Kubernetes manifests for doing all of that are not in this repository.

The easiest way to achieve all of this functionality is to run megaboom with [Deis Workflow](https://deis.com/docs/workflow).

Assuming you have a Workflow cluster running, and the `deis` CLI tool installed and configured to communicate with that cluster, run these commands to install megaboom:

```console
deis create --no-remote megaboom
deis pull quay.io/arschles/megaboom:devel -a megaboom
```

After you have megaboom installed, simply make an HTTP `POST` request to it to start a job. The command below shows how to do this, but it assumes you have the following environment variables set up:

- `DEIS_ROUTER_IP` - the IP address of the `deis-router` service. You can get this by executing the `kubectl get svc deis-router --namespace=deis` command
- `ENDPOINT` - the endpoint for megaboom to make requests against. This is the IP address or DNS name that you'd like to test under load.
- `HTTP_METHOD` - the HTTP Method (i.e. `GET`, `POST`, etc...) the execute against `ENDPOINT`
- `NUM_PODS` - the total number of pods to run in the Kubernetes job
- `NUM_CONCURRENT_PER_POD` - the number of concurrent requests to run in a single pod. The global maximum number of requests in flight at once will be `NUM_PODS` * `NUM_CONCURRENT_PER_POD`
- `NUM_REQUESTS_PER_POD` - the number of total requests to run in a single pod. The total number of requests to run globally will be `NUM_PODS` * `NUM_REQUESTS_PER_POD`
- `IMAGE` - the image to use for the job. Set this to `quay.io/arschles/megaboom:devel` (future versions of megaboom may support custom images)
- `JOB_NAMESPACE` - the Kubernetes namespace to run the job in

After you've set up all your environment variables, execute the following command:

```console
curl -XPOST -d '{"num_pods": ${NUM_PODS}, "num_concurrent_per_pod": ${NUM_CONCURRENT_PER_POD}, "num_reqs_per_pod": ${NUM_REQUESTS_PER_POD}, "http_method": "${HTTP_METHOD}", "endpoint": "${ENDPOINT}", "namespace": "${JOB_NAMESPACE}", "image": "${IMAGE}"}' http://megaboom.${DEIS_ROUTER_IP}.xip.io
```

Alternatively, you can look to the `Makefile` in this repository for a build target called `make-live`. This target conveniently can execute the above `curl` command given environment variables similar to the ones above.

However, the names of these variables are slightly different and the command assumes the HTTP server under test is also running in the same Deis Workflow cluster.
