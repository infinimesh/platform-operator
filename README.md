# Infinimesh Kubernetes Operator

Reconciles a Platform into the [core](https://github.com/infinimesh/core) stack: ArangoDB, RabbitMQ, Redis, and the core images, including `http-fs` on port 8000 with a volume at `/static`. On reconcile it deletes the old Kafka and Dgraph workloads (`device-registry`, `nodeserver`, `shadow-api`, `dgraph-alpha`, `dgraph-zero`, and the rest).

The static install of the same stack is [k8s](https://github.com/infinimesh/k8s).

A tag `vX.Y.Z` or `vX.Y.Z-rN` publishes `ghcr.io/infinimesh/platform-operator` under that tag and `:latest`. A revision tag also publishes `:vX.Y.Z`. A tag `vX.Y.Z-rc.N` publishes only that tag and is a GitHub pre-release. RC images older than 30 days are deleted.

## Core version

`spec.version` is the tag for every core service. An empty value is `v3.4.3`. `auto` checks the core GitHub releases once a day and uses the highest `vX.Y.Z` or `vX.Y.Z-rN`. It ignores `vX.Y.Z-rc.N`. `status.version` is the tag the stack is reconciled to.

```yaml
apiVersion: infinimesh.infinimesh.io/v1beta1
kind: Platform
metadata:
  name: my-infinimesh
spec:
  version: v3.4.3
```

Move a sticky stack to another tag:

```bash
kubectl patch platform my-infinimesh --type merge -p '{"spec":{"version":"v3.4.0"}}'
```

An RC is the same patch with `v3.4.0-rc.2`. ArangoDB, RabbitMQ, and Redis stay on their own images.

## Operator image

The controller runs as the StatefulSet `infinimesh-controller-manager` in `infinimesh-system`. It does not change its own image. Set the manager container, then apply the patch:

```bash
kubectl patch statefulset infinimesh-controller-manager -n infinimesh-system --patch-file manifests/operator-image.yaml
```

`manifests/operator-image.yaml` pins `ghcr.io/infinimesh/platform-operator:v3.4.3`. The same change with kubectl:

```bash
kubectl set image statefulset/infinimesh-controller-manager -n infinimesh-system manager=ghcr.io/infinimesh/platform-operator:v3.4.3
```

## Checks

`make commit-check` runs `fmt`, `race`, `test-fuzz`, and `code-ql-gate`. `fmt` writes Go files that `gofmt` would change, skipping `vendor`. `race` is `go test -race` for `./pkg/...` and `./cmd/...` with `-mod=vendor`. `code-ql` runs the CodeQL CLI and writes SARIF under `.tmp/codeql/`. The gate fails when a result is level `error`, ignoring paths under `vendor`. `jq` is required. GitHub Actions runs the same race and fuzz steps, and CodeQL on pull requests.

## Installation
```
kubectl apply -f manifests
```

or
```
kubectl apply -f https://raw.githubusercontent.com/infinimesh/platform-operator/main/manifests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/infinimesh/platform-operator/main/manifests/operator.yaml
```

Module `github.com/infinimesh/enterprise-operator`. Builds with Go 1.27. The Platform CRD is `apiextensions.k8s.io/v1`.

Apache-2.0.
