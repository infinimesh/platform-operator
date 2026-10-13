package platform

import (
	"os"
	"reflect"
	"strings"
	"testing"

	infinimeshv1beta1 "github.com/infinimesh/enterprise-operator/pkg/apis/infinimesh/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
)

func TestCoreImage(t *testing.T) {
	if got := coreImage("api", "v3.2.0"); got != "ghcr.io/infinimesh/core/api:v3.2.0" {
		t.Fatalf("core image %s", got)
	}
	if got := coreImage("api", ""); got != "ghcr.io/infinimesh/core/api:"+defaultCoreVersion {
		t.Fatalf("default core image %s", got)
	}
	if imagePullPolicy(coreImage("api", "v3.2.0")) != "IfNotPresent" {
		t.Fatal("a pinned core tag must not pull on every restart")
	}
	if imagePullPolicy("arangodb/arangodb:latest") != "Always" {
		t.Fatal("other images stay on Always")
	}
}

func TestRulesAndCoapAreInTheStack(t *testing.T) {
	body, err := os.ReadFile("stack.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, `suffix: "rules"`) || !strings.Contains(text, `suffix: "coap"`) || !strings.Contains(text, `suffix: "http-fs"`) {
		t.Fatal("rules, coap, or http-fs is missing from the stack")
	}
	if !strings.Contains(text, `coreImage("http-fs", instance.Status.Version)`) || !strings.Contains(text, `data: "/static"`) {
		t.Fatal("http-fs image or volume is missing")
	}
	if strings.Contains(text, "ports: []int{1883, 8883, 2112}") {
		t.Fatal("mqtt-bridge still publishes plaintext 1883")
	}
}

func TestLegacyDeleteIsGone(t *testing.T) {
	body, err := os.ReadFile("stack.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	for _, name := range []string{"deleteLegacy", "device-registry", "dgraph-zero", "dgraph-alpha"} {
		if strings.Contains(text, name) {
			t.Fatalf("stack.go still mentions %s", name)
		}
	}
}

func TestPlatformSpecHasNoKafka(t *testing.T) {
	if _, ok := reflect.TypeOf(infinimeshv1beta1.PlatformSpec{}).FieldByName("Kafka"); ok {
		t.Fatal("PlatformSpec still has Kafka")
	}
}

func TestIngressTLSIsNetworkingV1(t *testing.T) {
	field, ok := reflect.TypeOf(infinimeshv1beta1.PlatformApp{}).FieldByName("TLS")
	if !ok {
		t.Fatal("PlatformApp.TLS missing")
	}
	if field.Type.Elem() != reflect.TypeOf(networkingv1.IngressTLS{}) {
		t.Fatalf("PlatformApp.TLS is %s", field.Type)
	}
}

func TestInstallManifestsUseCurrentAPIs(t *testing.T) {
	for _, path := range []string{
		"../../../manifests/crd.yaml",
		"../../../config/crds/infinimesh_v1beta1_platform.yaml",
	} {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := string(body)
		if !strings.Contains(text, "apiVersion: apiextensions.k8s.io/v1\n") {
			t.Fatalf("%s is not an apiextensions v1 CRD", path)
		}
		if strings.Contains(text, "apiextensions.k8s.io/v1beta1") {
			t.Fatalf("%s still uses apiextensions v1beta1", path)
		}
	}
	for _, path := range []string{
		"../../../config/rbac/rbac_role.yaml",
		"../../../manifests/operator.yaml",
	} {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		text := string(body)
		if strings.Contains(text, "\n  - extensions\n") {
			t.Fatalf("%s still grants extensions ingress", path)
		}
		if !strings.Contains(text, "networking.k8s.io") || !strings.Contains(text, "persistentvolumeclaims") {
			t.Fatalf("%s is missing networking ingress or persistentvolumeclaims", path)
		}
	}
}
