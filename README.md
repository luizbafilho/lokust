# lokust

Locust.io Kubernetes Operator

## Initializing lokust on the Kubernetes cluster

This deploys all necessary k8s CRDs used to manage the tests

```sh
linkerd install | kubectl apply -f -
```

## Creating a new load test

**Creates the test loadtests namespace using a single file**

```sh
lokustctl create --namespace loadtests --name app-test --replicas 8 -f locustfile.py
```

**Creates the the test in the current namespace passing the directory as execution module**

> requires the locustfile.py file to be in the root directory

```sh
lokustctl create --name app-test --replicas 8 -f ./load-test
```

**Lists all current tests**

```sh
lokustctl ls
```

**Connect to a test**

It creates a kubernetes `port-forward` to the master pod, so you can access the dashboard and execute your test.

```sh
lokustctl connect --name app-test
```

**Delete a test**

```sh
lokustctl delete --name app-test
```
