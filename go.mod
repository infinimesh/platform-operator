module github.com/infinimesh/enterprise-operator

go 1.16

require (
	github.com/dgraph-io/dgo/v210 v210.0.0-20210421093152-78a2fece3ebd
	github.com/go-logr/logr v0.4.0
	github.com/infinimesh/infinimesh v1.0.0
	github.com/onsi/gomega v1.13.0
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	google.golang.org/grpc v1.31.1
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
	k8s.io/code-generator v0.21.2 // indirect
	sigs.k8s.io/apiserver-network-proxy/konnectivity-client v0.0.19 // indirect
	sigs.k8s.io/controller-runtime v0.9.0
)
