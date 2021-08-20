# megaboom

Megaboom is a Distributed HTTP Load Generator, based on Kubernetes. It is an HTTP server meant to be run inside your Kubernetes cluster, behind a [Service](http://kubernetes.io/docs/user-guide/services/), and it runs [Jobs](http://kubernetes.io/docs/user-guide/jobs/) that invocate the [hey](https://github.com/rakyll/hey) CLI (formerly called `boom`, hence the name of this repository) to generate load.

## Alpha Status

This project is still in alpha, and is missing a few major features:

- **Reporting the success/failure rate of ongoing jobs** - currently, it's expected that you introspect each [Pod](http://kubernetes.io/docs/user-guide/pods/) to get individual `boom` command output, or use your own monitoring/metrics system to see status of your application under load
- **Cleaning up after itself** - currently it creates a `megaboom` job in the `default` namespace. You must delete that job before you make the next HTTP request to megaboom. Do the deletion with `kubectl delete job megaboom`
- **CI infrastructure to run tests and automatically build images** - currently, there are very few tests, and the current Docker image has not been pushed anywhere. You have to build it yourself and push to your own registry.

## Usage

This server must run inside a Kubernetes pod. It's recommended that you run this server in a pod managed by a `Deployment`, and then run a service in front of the replication controller. There is a [Helm](https://helm.sh) chart available for you in this repository. Install it with the following command:

```shell
helm install megaboom -n megaboom ./chart/megaboom
```

## Making Requests to Megaboom

Megaboom starts a (simple) HTTP server, which you make requests to via an HTTP client (e.g. `curl`). The helm command in the previous section installs a `megaboom` service in the `megaboom` namespace, and installs a separate "debug pod" to which you can log in and make these requests.

>You can also issue requests to this server using Kubernetes port-forwarding or proxying if you'd like.

>The debug pod does not have `curl` installed by default. The first time you open a shell on it, execute `apt-get update && apt-get install -y curl`. As long as that same pod exists thereafter, you won't need to reinstall `curl`.

Open a shell on the debug pod and run the following `curl` command:

```shell
curl -d '{"num_runners": 50, "num_concurrent_per_runner": 2000, "num_reqs_per_runner": 200000, "endpoint": "https://gifm.dev", "namespace": "megaboom"}' http://megaboom:8080/job
```

This request will return a Job ID. Copy this down, and use it to delete the job later:

```shell
curl -XDELETE http://megaboom:8080/job/<job_id>
```
