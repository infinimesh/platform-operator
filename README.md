# Infinimesh Kubernetes Operator

Reconciles a Platform into the [core](https://github.com/infinimesh/core) stack: ArangoDB, RabbitMQ, Redis, and the images `ghcr.io/infinimesh/core/<service>:latest`, including `http-fs` on port 8000 with a volume at `/static`. On reconcile it deletes the old Kafka and Dgraph workloads (`device-registry`, `nodeserver`, `shadow-api`, `dgraph-alpha`, `dgraph-zero`, and the rest).

The static install of the same stack is [k8s](https://github.com/infinimesh/k8s).

A tag `vX.Y.Z` or `vX.Y.Z-rN` publishes `ghcr.io/infinimesh/platform-operator` under that tag and `:latest`. A revision tag also publishes `:vX.Y.Z`.

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
