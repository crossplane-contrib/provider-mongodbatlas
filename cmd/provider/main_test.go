package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestProviderGroups(t *testing.T) {
	s := runtime.NewScheme()
	for _, g := range []string{
		"mongodbatlas.crossplane.io",
		"cloud.mongodbatlas.crossplane.io",
		"cloud.mongodbatlas.m.crossplane.io",
		"apiextensions.k8s.io",
		"notmongodbatlas.crossplane.io",
	} {
		s.AddKnownTypeWithName(schema.GroupVersionKind{Group: g, Version: "v1", Kind: "Secret"}, &corev1.Secret{})
	}

	got := providerGroups(s, "mongodbatlas.crossplane.io", "mongodbatlas.m.crossplane.io")
	want := []string{"cloud.mongodbatlas.crossplane.io", "cloud.mongodbatlas.m.crossplane.io", "mongodbatlas.crossplane.io"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("providerGroups(...): -want, +got:\n%s", diff)
	}
}
