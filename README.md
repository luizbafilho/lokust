# lokust

Locust.io Kubernetes Operator

## Initializing lokust on the Kubernetes cluster

This deploys all necessary k8s CRDs used to manage the tests

```sh
lokustctl init
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

**Delete a test**

```sh
lokustctl delete --name app-test-a8djvhi
```
