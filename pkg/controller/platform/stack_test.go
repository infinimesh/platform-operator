package platform

import (
	"reflect"
	"testing"

	infinimeshv1beta1 "github.com/infinimesh/enterprise-operator/pkg/apis/infinimesh/v1beta1"
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
