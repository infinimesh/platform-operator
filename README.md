# Infinimesh Kubernetes Operator

Reconciles a Platform into the [core](https://github.com/infinimesh/core) stack: ArangoDB, RabbitMQ, Redis, and the images `ghcr.io/infinimesh/core/<service>:latest`. On reconcile it deletes the old Kafka and Dgraph workloads (`device-registry`, `nodeserver`, `shadow-api`, `dgraph-alpha`, `dgraph-zero`, and the rest).

The static install of the same stack is [k8s](https://github.com/infinimesh/k8s).

A tag `vX.Y.Z` or `vX.Y.Z-rN` publishes `ghcr.io/infinimesh/platform-operator` under that tag and `:latest`. A revision tag also publishes `:vX.Y.Z`.

## Installation
```
kubectl apply -f manifests
```

or
```
kubectl apply -f https://raw.githubusercontent.com/infinimesh/platform-operator/main/manifests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/infinimesh/platform-operator/main/manifests/operator.yaml
```

Module `github.com/infinimesh/enterprise-operator`.

Apache-2.0.
