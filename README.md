# lokust

Lokust is a [Locust.io](https://locust.io/) Kubernetes Operator. It helps integrate Locust distributed mode into Kubernetes, so you easely scale up and down your load tests without the hassle of provisioning infrastructure.

## Getting Started

**Install lokustctl**

It will be used to interface with the `lokust` controller, helping to manage your tests without having to deal with kubernetes yaml files, so anyone regardless of knowing kubernetes can create and run tests.

```sh
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
VERSION="0.1.0-beta.1"

curl -L -s -o ./lokustctl "https://github.com/luizbafilho/lokust/releases/download/v${VERSION}/lokustctl_${VERSION}_${OS}_${ARCH}"

chmod +x ./lokustctl
mv lokustctl /usr/local/bin
lokustctl -h
```

**Deploy Lokust Controller**

You can you `lokustctl` to deploy the Lokust controller into Kubernetes. The following command will apply all kubernetes manifests necessary  deploy the controller.

It will deploy them by default into `lokust-system` namespace, you can specify a different one using `--namespace`

```sh
$ lokustctl install | kubectl apply -f -
```


**Create a new load test**

We first need a `locustfile.py` that defines your test behavior. More info at [locust.io](https://locust.io/) on how to write it.

```sh
cat <<EOT >> locustfile.py
from locust import HttpUser, between, task

class WebsiteUser(HttpUser):
    host = "https://google.com"
    wait_time = between(5, 15)

    @task
    def index(self):
        self.client.get("/")
EOT
```

```sh
lokustctl create --name app-test --replicas 2 -f locustfile.py
```

**Lists all current tests**

```sh
lokustctl list
```

**Access the test's dashboard**

To actually start the test you need to access the test's dashboard. To access it in a kubernetes environment we need to create a local proxy that redirects you the the test created there. To do that execute:

```sh
lokustctl connect app-test
```

**Delete the test**

Due the locust tests nature, it will run indefinetly on kubernetes consuming resources until you delete it.

```sh
lokustctl delete app-test
```
