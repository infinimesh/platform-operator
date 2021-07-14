module github.com/infinimesh/enterprise-operator

go 1.15

require (
	cloud.google.com/go v0.54.0 // indirect
	github.com/appscode/jsonpatch v1.0.1 // indirect
	github.com/dgraph-io/dgo/v210 v210.0.0-20210421093152-78a2fece3ebd
	github.com/evanphx/json-patch v4.11.0+incompatible // indirect
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/infinimesh/infinimesh v1.0.0
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/prometheus/client_golang v1.11.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/grpc v1.31.1
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver v0.21.2 // indirect
	k8s.io/apimachinery v0.21.2
	k8s.io/client-go v0.21.2
	k8s.io/klog v0.3.1 // indirect
	k8s.io/kube-openapi v0.0.0-20210527164424-3c818078ee3d // indirect
	sigs.k8s.io/controller-runtime v0.2.0-alpha.0
	sigs.k8s.io/testing_frameworks v0.1.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/api => k8s.io/api v0.0.0-20181213150558-05914d821849
replace sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.10
replace k8s.io/client-go => k8s.io/client-go v0.0.0-20181213151034-8d9ed539ba31
replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93
replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20181213153335-0fe22c71c476
replace github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.2.0
replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190306001800-15615b16d372
replace k8s.io/klog => k8s.io/klog v0.2.0
replace sigs.k8s.io/testing_frameworks => sigs.k8s.io/testing_frameworks v0.1.1
replace sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.2.0