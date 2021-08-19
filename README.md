# megaboom

Megaboom is a Distributed HTTP Load Generator, based on Kubernetes. It is an HTTP server meant to be run inside your Kubernetes cluster, behind a [Service](http://kubernetes.io/docs/user-guide/services/), and it runs [Jobs](http://kubernetes.io/docs/user-guide/jobs/) that invocate the [hey](https://github.com/rakyll/hey) CLI (formerly called `boom`, hence the name of this repository) to generate load.

## Alpha Status

This project is still in alpha, and is missing a few major features:

- **Reporting the success/failure rate of ongoing jobs** - currently, it's expected that you introspect each [Pod](http://kubernetes.io/docs/user-guide/pods/) to get individual `boom` command output, or use your own monitoring/metrics system to see status of your application under load
- **Cleaning up after itself** - currently it creates a `megaboom` job in the `default` namespace. You must delete that job before you make the next HTTP request to megaboom. Do the deletion with `kubectl delete job megaboom`
- **CI infrastructure to run tests and automatically build images** - currently, there are very few tests, and the current Docker image has been manually pushed to `quay.io/arschles/megaboom:devel`

## Usage

This server must run inside a Kubernetes pod. It's recommended that you run this server in a pod managed by a `Deployment`, and then run a service in front of the replication controller. There is a [Helm](https://helm.sh) chart available for you in this repository. Install it with the following command:

```shell
helm install megaboom -n megaboom ./chart/megaboom
```

## Making Requests to Megaboom

After you have megaboom installed, simply make an HTTP `POST` request to it to start a job. The command below shows how to do this, but it assumes you have the following environment variables set up:

- `DEIS_ROUTER_IP` - the IP address of the `deis-router` service. You can get this by executing the `kubectl get svc deis-router --namespace=deis` command
- `ENDPOINT` - the endpoint for megaboom to make requests against. This is the IP address or DNS name that you'd like to test under load.
- `TEST_NUM_PODS` - the total number of pods to run in the Kubernetes job
- `TEST_NUM_CONCURRENT_PER_POD` - the number of concurrent requests to run in a single pod. The global maximum number of requests in flight at once will be `NUM_PODS` * `NUM_CONCURRENT_PER_POD`
- `TEST_NUM_REQS_PER_POD` - the number of total requests to run in a single pod. The total number of requests to run globally will be `NUM_PODS` * `NUM_REQUESTS_PER_POD`

After you've set up all your environment variables, execute the following command:

```console
curl -XPOST -d '{"num_pods": ${TEST_NUM_PODS}, "num_concurrent_per_pod": ${TEST_NUM_CONCURRENT_PER_POD}, "num_reqs_per_pod": ${TEST_NUM_REQS_PER_POD}, "http_method": "GET", "endpoint": "${ENDPOINT}", "namespace": "default", "image": "quay.io/arschles/megaboom:devel"}' http://megaboom.${DEIS_ROUTER_IP}.xip.io
```

## Makefile Convenience Target

Alternatively, there is a build target in the `Makefile` in this repository called `make-live`. This target conveniently can execute the above `curl` command given environment variables similar to the ones above. This target assumes the HTTP server under test is also running in the same Deis Workflow cluster.

Below are the environment variables to use:

- `DEIS_ROUTER_IP` - the IP address of the `deis-router` service. You can get this by executing the `kubectl get svc deis-router --namespace=deis` command
- `APP_NAME` - the name of the Deis Workflow app to test
- `TEST_NUM_PODS` - the total number of pods to run in the Kubernetes job
- `TEST_NUM_CONCURRENT_PER_POD` - the number of concurrent requests to run in a single pod. The global maximum number of requests in flight at once will be `NUM_PODS` * `NUM_CONCURRENT_PER_POD`
- `TEST_NUM_REQS_PER_POD` - the number of total requests to run in a single pod. The total number of requests to run globally will be `NUM_PODS` * `NUM_REQUESTS_PER_POD`

After you've configured all of these environment variables, simply execute the below command to start the test:

```console
make test-live
```
