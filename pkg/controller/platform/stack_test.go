package platform

import (
	"reflect"
	"testing"

	infinimeshv1beta1 "github.com/infinimesh/enterprise-operator/pkg/apis/infinimesh/v1beta1"
)

func TestCoreImage(t *testing.T) {
	got := coreImage("repo")
	want := "ghcr.io/infinimesh/core/repo:latest"
	if got != want {
		t.Fatalf("core image %s", got)
	}
}

func TestPlatformSpecHasNoKafka(t *testing.T) {
	if _, ok := reflect.TypeOf(infinimeshv1beta1.PlatformSpec{}).FieldByName("Kafka"); ok {
		t.Fatal("PlatformSpec still has Kafka")
	}
}
