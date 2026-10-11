# Infinimesh Kubernetes Operator

[![Docker Repository on Quay](https://quay.io/repository/infinimesh/platform-operator/status "Docker Repository on Quay")](https://quay.io/repository/infinimesh/platform-operator)

Installs an Infinimesh platform into a Kubernetes cluster.

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
