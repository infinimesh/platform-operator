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
	got := coreImage("api")
	want := "ghcr.io/infinimesh/core/api:latest"
	if got != want {
		t.Fatalf("core image %s", got)
	}
}

func TestLegacyNamesAreRemoved(t *testing.T) {
	want := map[string]bool{"repo": false, "web": false, "handsfree": false}
	for _, suffix := range legacySuffixes {
		if _, ok := want[suffix]; ok {
			want[suffix] = true
		}
	}
	for name, ok := range want {
		if !ok {
			t.Fatalf("legacy suffix %s is not deleted", name)
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
