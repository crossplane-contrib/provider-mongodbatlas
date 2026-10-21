/*
Copyright 2026 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/crossplane/crossplane-runtime/v2/pkg/logging"
	diffv1alpha1 "github.com/crossplane/upjet/v2/proto/diff/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/crossplane-contrib/provider-mongodbatlas/config"
)

// The manifests below are what a client puts in a plan request. Every
// ProviderConfig references a Secret that does not exist anywhere, which the
// offline diff server must never need.
const (
	clusterProviderConfig = `
apiVersion: mongodbatlas.m.crossplane.io/v1beta1
kind: ClusterProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: provider-credentials
      namespace: crossplane-system
      key: credentials
`
	legacyProviderConfig = `
apiVersion: mongodbatlas.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: provider-credentials
      namespace: crossplane-system
      key: credentials
`
	// A namespaced Project, a plugin framework resource, in a namespace of its
	// own, referencing the cluster-scoped ProviderConfig above.
	namespacedProject = `
apiVersion: mongodbatlas.m.crossplane.io/v1alpha1
kind: Project
metadata:
  name: example
  namespace: team-a
spec:
  forProvider:
    name: example-project
    orgId: 5e0f6ad4f4e9f9b1c2d3e4a1
  providerConfigRef:
    kind: ClusterProviderConfig
    name: default
`
	// The actual state of an Auditing, a plugin SDKv2 resource, as the
	// provider records it after a reconcile: the external name annotation and
	// the observation the plan is computed against.
	clusterAuditingActual = `
apiVersion: mongodbatlas.crossplane.io/v1alpha1
kind: Auditing
metadata:
  name: example
  annotations:
    crossplane.io/external-name: 5e0f6ad4f4e9f9b1c2d3e4f5
spec:
  forProvider:
    auditAuthorizationSuccess: false
    auditFilter: '{"atype": "authenticate"}'
    enabled: true
    projectId: 5e0f6ad4f4e9f9b1c2d3e4f5
  providerConfigRef:
    name: default
status:
  atProvider:
    id: 5e0f6ad4f4e9f9b1c2d3e4f5
    auditAuthorizationSuccess: false
    auditFilter: '{"atype": "authenticate"}'
    configurationType: FILTER_JSON
    enabled: true
    projectId: 5e0f6ad4f4e9f9b1c2d3e4f5
`
	clusterAuditingDisabled = `
apiVersion: mongodbatlas.crossplane.io/v1alpha1
kind: Auditing
metadata:
  name: example
spec:
  forProvider:
    auditAuthorizationSuccess: false
    auditFilter: '{"atype": "authenticate"}'
    enabled: false
    projectId: 5e0f6ad4f4e9f9b1c2d3e4f5
  providerConfigRef:
    name: default
`
	clusterAuditingMoved = `
apiVersion: mongodbatlas.crossplane.io/v1alpha1
kind: Auditing
metadata:
  name: example
spec:
  forProvider:
    auditAuthorizationSuccess: false
    auditFilter: '{"atype": "authenticate"}'
    enabled: true
    projectId: 5e0f6ad4f4e9f9b1c2d3e4f6
  providerConfigRef:
    name: default
`
	// A ResourcePolicy: the upstream resource validates its policies against
	// the Atlas API while planning, which the offline server refuses.
	namespacedResourcePolicy = `
apiVersion: mongodbatlas.m.crossplane.io/v1alpha1
kind: ResourcePolicy
metadata:
  name: forbid-access-from-anywhere
  namespace: team-a
spec:
  forProvider:
    name: forbid-access-from-anywhere
    orgId: 5e0f6ad4f4e9f9b1c2d3e4a1
    policies:
    - body: |
        forbid (principal, action == ResourcePolicy::Action::"project.ipAccessList.modify", resource)
        when { context.project.ipAccessList.contains(ip("0.0.0.0/0")) };
  providerConfigRef:
    kind: ClusterProviderConfig
    name: default
`
	missingProviderConfigProject = `
apiVersion: mongodbatlas.m.crossplane.io/v1alpha1
kind: Project
metadata:
  name: example
  namespace: team-a
spec:
  forProvider:
    name: example-project
    orgId: 5e0f6ad4f4e9f9b1c2d3e4a1
  providerConfigRef:
    kind: ProviderConfig
    name: missing
`
)

// TestDiffServerPlan drives the diff gRPC server the way a client does, over
// a unix socket, with the provider configurations and the offline Terraform
// setup main wires it with. It covers a namespaced plugin framework resource
// resolving a cluster-scoped ProviderConfig from its own namespace, every plan
// action on a plugin SDKv2 resource, and a resource whose plan needs the
// Atlas API.
func TestDiffServerPlan(t *testing.T) {
	ctx := context.Background()
	s, scheme, err := newDiffServer(logging.NewNopLogger(), config.GetidentifierFromProvider(), config.GetProviderNamespaced())
	if err != nil {
		t.Fatalf("cannot initialize the diff gRPC server: %v", err)
	}

	// A unix socket, because Serve does not report the port it bound when
	// given port zero. Its path is kept short: macOS limits socket paths to
	// 104 bytes, which the path t.TempDir returns can exceed.
	dir, err := os.MkdirTemp("", "diff")
	if err != nil {
		t.Fatalf("cannot create the socket directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) }) //nolint:errcheck // best effort cleanup of a temporary directory
	sock := filepath.Join(dir, "diff.sock")

	serveCtx, cancel := context.WithCancel(ctx)
	served := make(chan error, 1)
	go func() { served <- s.Serve(serveCtx, "unix", sock, scheme) }()
	t.Cleanup(func() {
		cancel()
		if err := <-served; err != nil {
			t.Errorf("the diff gRPC server failed: %v", err)
		}
	})

	conn, err := grpc.NewClient("unix://"+sock, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("cannot connect to the diff gRPC server: %v", err)
	}
	t.Cleanup(func() { conn.Close() }) //nolint:errcheck // the process is exiting
	client := diffv1alpha1.NewPlanServiceClient(conn)

	type want struct {
		action  diffv1alpha1.Action
		changes map[string]string
		code    codes.Code
		err     string
	}
	cases := map[string]struct {
		desired string
		actual  string
		store   []string
		want    want
	}{
		"NamespacedProjectCreateWithClusterProviderConfig": {
			desired: namespacedProject,
			store:   []string{clusterProviderConfig},
			want: want{
				action: diffv1alpha1.Action_ACTION_CREATE,
				changes: map[string]string{
					"spec.forProvider.name":  "<unset> -> example-project",
					"spec.forProvider.orgId": "<unset> -> 5e0f6ad4f4e9f9b1c2d3e4a1",
				},
			},
		},
		"ClusterAuditingUpdate": {
			desired: clusterAuditingDisabled,
			actual:  clusterAuditingActual,
			store:   []string{legacyProviderConfig},
			want: want{
				action:  diffv1alpha1.Action_ACTION_UPDATE,
				changes: map[string]string{"spec.forProvider.enabled": "true -> false"},
			},
		},
		"ClusterAuditingReplace": {
			desired: clusterAuditingMoved,
			actual:  clusterAuditingActual,
			store:   []string{legacyProviderConfig},
			want: want{
				action:  diffv1alpha1.Action_ACTION_REPLACE,
				changes: map[string]string{"spec.forProvider.projectId": "5e0f6ad4f4e9f9b1c2d3e4f5 -> 5e0f6ad4f4e9f9b1c2d3e4f6 (requires replace)"},
			},
		},
		"ClusterAuditingNoOp": {
			desired: clusterAuditingActual,
			actual:  clusterAuditingActual,
			store:   []string{legacyProviderConfig},
			want:    want{action: diffv1alpha1.Action_ACTION_NO_OP},
		},
		"NamespacedResourcePolicyNeedsTheAPI": {
			desired: namespacedResourcePolicy,
			store:   []string{clusterProviderConfig},
			want:    want{code: codes.FailedPrecondition, err: "diff computation is not supported"},
		},
		"MissingProviderConfig": {
			desired: missingProviderConfigProject,
			store:   []string{clusterProviderConfig},
			want:    want{code: codes.Internal, err: "cannot get referenced ProviderConfig"},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			req := &diffv1alpha1.PlanRequest{
				DesiredResource: asStruct(t, tc.desired),
				ActualResource:  asStruct(t, tc.actual),
			}
			for _, m := range tc.store {
				req.KubernetesObjectStore = append(req.KubernetesObjectStore, asStruct(t, m))
			}

			callCtx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()
			rsp, err := client.Plan(callCtx, req, grpc.WaitForReady(true))
			if tc.want.err != "" {
				if err == nil {
					t.Fatalf("expected the plan to fail with %q, got response: %s", tc.want.err, protojson.Format(rsp))
				}
				if st := status.Convert(err); st.Code() != tc.want.code || !strings.Contains(st.Message(), tc.want.err) {
					t.Fatalf("expected a %s error containing %q, got %s: %s", tc.want.code, tc.want.err, st.Code(), st.Message())
				}
				return
			}
			if err != nil {
				t.Fatalf("cannot plan: %v", err)
			}
			if diff := cmp.Diff(tc.want.action.String(), rsp.GetAction().String()); diff != "" {
				t.Errorf("unexpected plan action, -want +got:\n%s", diff)
			}
			if diff := cmp.Diff(tc.want.changes, declaredChanges(rsp)); diff != "" {
				t.Errorf("unexpected changes to the declared fields, -want +got:\n%s\nfull response: %s", diff, protojson.Format(rsp))
			}
		})
	}
}

// collectionIndex matches the index of a collection element in a field path.
// A set element is indexed by the hash of its content, which no expectation
// should spell out, so indices are rendered as [*] instead.
var collectionIndex = regexp.MustCompile(`\[\d+\]`)

// declaredChanges renders the changes a plan reports for the fields the
// desired resource declares, as field to "actual -> planned". The changes the
// Terraform provider contributes on its own, such as the defaults it fills in
// on a create, are left out, so that a provider bump changing a default does
// not turn into a failing plan here.
func declaredChanges(rsp *diffv1alpha1.PlanResponse) map[string]string {
	changes := map[string]string{}
	for _, c := range rsp.GetChanges() {
		if c.GetOrigin() != diffv1alpha1.Origin_ORIGIN_DESIRED_STATE {
			continue
		}
		change := renderValue(c.GetActual()) + " -> " + renderValue(c.GetPlanned())
		if c.GetRequiresReplace() {
			change += " (requires replace)"
		}
		changes[collectionIndex.ReplaceAllString(c.GetField(), "[*]")] = change
	}
	if len(changes) == 0 {
		return nil
	}
	return changes
}

func renderValue(v *diffv1alpha1.FieldValue) string {
	switch {
	case v == nil:
		return "<unset>"
	case v.GetValue() != nil:
		return fmt.Sprintf("%v", v.GetValue().AsInterface())
	default:
		return v.GetAbsence().String()
	}
}

// asStruct converts a YAML manifest into the protobuf Struct a plan request
// carries it as. An empty manifest converts to nil, which is how a request
// says that the resource does not exist yet.
func asStruct(t *testing.T, manifest string) *structpb.Struct {
	t.Helper()
	if manifest == "" {
		return nil
	}
	j, err := yaml.ToJSON([]byte(manifest))
	if err != nil {
		t.Fatalf("cannot convert the manifest to JSON: %v", err)
	}
	s := &structpb.Struct{}
	if err := s.UnmarshalJSON(j); err != nil {
		t.Fatalf("cannot convert the manifest to a protobuf Struct: %v", err)
	}
	return s
}
