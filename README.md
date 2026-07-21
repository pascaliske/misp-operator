# MISP-Operator

> A Kubernetes operator for simplified deployments of MISP at scale.

This project aims to simplify the deployment, configuration and management of any number of MISP instances inside a Kubernetes cluster through a set of CRDs.
It builds upon the images from the [misp/misp-docker](https://github.com/misp/misp-docker) project and integrates them with native Kubernetes tooling.

> [!NOTE]
> This project is still in the early stages of development. Feedback is highly appreciated!

## Features

- Deploying MISP instances easily via [`MispInstance`](config/samples/mispinstance-minimal.yaml)
- Built-in deployment of [MISP modules](https://github.com/misp/misp-modules) container
- Seamless integration with native K8s tooling, e.g. [External Secrets Operator](https://github.com/external-secrets/external-secrets) or [cert-manager](https://github.com/cert-manager/cert-manager)
- GitOps friendly by design

## Getting Started

### Installation

#### Option 1: Install via Helm Chart

There is a Helm Chart available at [ghcr.io/pascaliske/charts/misp-operator](https://ghcr.io/pascaliske/charts/misp-operator).

You can install it using Helm as follows:

```shell
helm upgrade --install --namespace misp-operator-system --create-namespace misp-operator oci://ghcr.io/pascaliske/charts/misp-operator
```

#### Option 2: Install via Kubectl

There is also an install manifest attached to all releases which can be used to install the operator via `kubectl`:

```shell
kubectl apply -n misp-operator-system -f https://github.com/pascaliske/misp-operator/releases/latest/download/install.yaml
```

### Usage

As minimal example the following `MispInstance` can be used:

```yaml
apiVersion: misp.k8s.pascaliske.dev/v1alpha1
kind: MispInstance
metadata:
  name: misp-minimal
spec:
  # pause reconciliation of this resource
  suspend: false

  # provide the front-facing url of the instance
  baseUrl: https://misp.example.org

  # provide database host & credentials - port defaults to '3306', database defaults to 'misp'
  database:
    host: mariadb
    credentialsSecretRef:
      name: mariadb-user

  # provide redis/valkey host & password - port defaults to '6379'
  cache:
    host: redis
    passwordSecretRef:
      name: redis-user
```

You can find more examples in the [config/samples](config/samples/) directory.

## Contributing

Contributions of any kind are welcome! Please see the [contribution guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE.md) file for details.
