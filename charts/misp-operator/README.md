# [`misp-operator`](https://github.com/pascaliske/misp-operator)

> Helm chart for MISP-Operator - an operator for simplified deployments of MISP at scale

<!-- x-release-please-start-version -->
[![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)](https://github.com/pascaliske/misp-operator)
[![Version](https://img.shields.io/static/v1?label=Version&message=0.0.9&color=informational&style=flat-square)](https://github.com/pascaliske/misp-operator)
[![AppVersion](https://img.shields.io/static/v1?label=AppVersion&message=0.0.9&color=informational&style=flat-square)](https://github.com/pascaliske/misp-operator)
<!-- x-release-please-end -->

* <https://github.com/pascaliske/misp-operator>
* <https://github.com/misp/misp-docker>
* <https://github.com/misp/misp>

## Requirements

- [`helm`](https://helm.sh) - Refer to their [docs](https://helm.sh/docs) to get started.

## Usage

To install this chart simply run the following command:

```shell
helm upgrade --install --namespace misp-operator-system --create-namespace misp-operator oci://ghcr.io/pascaliske/charts/misp-operator
```

To uninstall this chart simply run the following command:

```sh
helm delete misp-operator
```

## Verification

The OCI chart of the operator is **keylessly** signed using [`cosign`](https://docs.sigstore.dev/cosign/verifying/verify/) and can be can be verified:

<!-- x-release-please-start-version -->
```shell
cosign verify ghcr.io/pascaliske/charts/misp-operator:0.0.9 \
  --certificate-identity-regexp "^https://github.com/pascaliske/misp-operator.*$" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```
<!-- x-release-please-end -->

> [!NOTE]
> Verification succeeds only if the artifact was signed by the GitHub Actions workflow in this repository.
> Any modification of the artifact or signing from a different identity will cause verification to fail.

## Values

The following values can be used to adjust the helm chart.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Pod-level affinity. More info [here](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#scheduling). |
| certManager.enabled | bool | `true` | Enable the cert-manager integration for providing certs for metrics and webhooks. |
| certManager.issuerRef | object | `{}` | Provide custom issuer for certs. |
| controller.annotations | object | `{}` | Additional annotations for the controller object. |
| controller.labels | object | `{}` | Additional labels for the controller object. |
| controller.replicas | int | `1` | The number of replicas. |
| controller.updateStrategy | object | `{}` | The controller update strategy. Currently only applies to controllers of kind `Deployment`. |
| crds.create | bool | `true` | Specifies whether the CRDs should be created when installing the chart. |
| extraArgs | list | `[]` | Additional arguments to be added to the operator's args list. |
| extraVolumeMounts | list | `[]` | Specify extra volume mounts for the default containers. |
| extraVolumes | list | `[]` | Specify extra volumes for the workload. |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` | The pull policy for the controller. |
| image.registry | string | `"ghcr.io"` | The registry to pull the image from. |
| image.repository | string | `"pascaliske/misp-operator"` | The repository to pull the image from. |
| image.tag | string | `.Chart.AppVersion` | The image tag, if left empty chart's appVersion will be used. |
| imagePullSecrets | list | `[]` |  |
| metrics.enabled | bool | `true` | Enable operator internal metrics. Prometheus must be installed in the cluster |
| metrics.secure | bool | `true` | Enable secure metrics endpoint. This requires cert-manager to be installed. |
| metrics.serviceMonitor.annotations | object | `{}` | Additional annotations for the service monitor object. |
| metrics.serviceMonitor.enabled | bool | `true` | Create a service monitor for Prometheus operator. |
| metrics.serviceMonitor.interval | string | `"30s"` | How frequently the exporter should be scraped. |
| metrics.serviceMonitor.labels | object | `{}` | Additional labels for the service monitor object. |
| metrics.serviceMonitor.timeout | string | `"10s"` | Timeout value for individual scrapes. |
| nameOverride | string | `""` |  |
| namespaceOverride | string | `""` |  |
| nodeSelector | object | `{}` | Pod-level node selector. More info [here](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#scheduling). |
| podAnnotations | object | `{}` | Annotations to be added to the pod. |
| podLabels | object | `{}` | Labels to be added to the pod. |
| podSecurityContext | object | `{"runAsGroup":65532,"runAsNonRoot":true,"runAsUser":65532,"seccompProfile":{"type":"RuntimeDefault"}}` | Pod-level security attributes. More info [here](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#security-context). |
| priorityClassName | string | `""` | Optional priority class name indicates the importance of a pod relative to other pods. |
| rbac.create | bool | `true` | Specifies whether ClusterRole and ClusterRoleBinding should be created. |
| resources | object | `{"limits":{"cpu":"250m","memory":"256Mi"},"requests":{"cpu":"100m","memory":"128Mi"}}` | Compute resources used by the container. More info [here](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/). |
| securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true}` | Container-level security attributes. More info [here](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#security-context). |
| serviceAccount.annotations | object | `{}` | Additional annotations for the role and role binding objects. |
| serviceAccount.automount | bool | `true` | Automatically mount the ServiceAccount API credentials |
| serviceAccount.create | bool | `true` | Create a `ServiceAccount` object. |
| serviceAccount.labels | object | `{}` | Additional labels for the role and role binding objects. |
| serviceAccount.name | string | `""` | Specify the service account used for the controller. |
| tolerations | list | `[]` | Pod-level tolerations. More info [here](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/pod-v1/#scheduling). |
| webhook.enabled | bool | `true` | Enable the webhooks endpoint. |

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| pascaliske | <info@pascaliske.dev> | <https://pascaliske.dev> |

## Contributing

Contributions of any kind are welcome! Please see the [contribution guide](../../CONTRIBUTING.md) for details.

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](../../LICENSE.md) file for details.
