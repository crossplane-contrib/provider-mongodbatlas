package refs

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestNotFoundWhenPathParamEmpty(t *testing.T) {
	errDiag := func(detail string) *tfprotov6.Diagnostic {
		return &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "error getting Alert Configuration information: %s",
			Detail:   detail,
		}
	}
	warning := &tfprotov6.Diagnostic{Severity: tfprotov6.DiagnosticSeverityWarning, Detail: "deprecated attribute"}
	cases := map[string]struct {
		diags []*tfprotov6.Diagnostic
		want  bool
	}{
		"NoDiagnostics":          {diags: nil, want: false},
		"OnlyWarnings":           {diags: []*tfprotov6.Diagnostic{warning}, want: false},
		"EmptyGroupID":           {diags: []*tfprotov6.Diagnostic{errDiag("groupId is empty and must be specified")}, want: true},
		"EmptyPathParamAndWarn":  {diags: []*tfprotov6.Diagnostic{warning, errDiag("alertConfigId is empty and must be specified")}, want: true},
		"APIError":               {diags: []*tfprotov6.Diagnostic{errDiag("https://cloud.mongodb.com/api/atlas/v2/groups/x/alertConfigs/y GET: HTTP 401 Unauthorized")}, want: false},
		"EmptyPathParamAndOther": {diags: []*tfprotov6.Diagnostic{errDiag("groupId is empty and must be specified"), errDiag("HTTP 500")}, want: false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := NotFoundWhenPathParamEmpty(tc.diags)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NotFoundWhenPathParamEmpty(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestNotFoundWhenErrorSummary(t *testing.T) {
	errDiag := func(summary string) *tfprotov6.Diagnostic {
		return &tfprotov6.Diagnostic{Severity: tfprotov6.DiagnosticSeverityError, Summary: summary, Detail: "detail"}
	}
	warning := &tfprotov6.Diagnostic{Severity: tfprotov6.DiagnosticSeverityWarning, Summary: "API key not found"}
	cases := map[string]struct {
		diags []*tfprotov6.Diagnostic
		want  bool
	}{
		"NoDiagnostics":       {diags: nil, want: false},
		"OnlyWarnings":        {diags: []*tfprotov6.Diagnostic{warning}, want: false},
		"MatchingSummary":     {diags: []*tfprotov6.Diagnostic{errDiag("API key not found")}, want: true},
		"MatchingAndWarning":  {diags: []*tfprotov6.Diagnostic{warning, errDiag("API key not found")}, want: true},
		"OtherSummary":        {diags: []*tfprotov6.Diagnostic{errDiag("error fetching resource")}, want: false},
		"MatchingAndOtherErr": {diags: []*tfprotov6.Diagnostic{errDiag("API key not found"), errDiag("error fetching resource")}, want: false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := NotFoundWhenErrorSummary("API key not found")(tc.diags)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("NotFoundWhenErrorSummary(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestSetIdentifierArgument(t *testing.T) {
	type args struct {
		base         map[string]any
		externalName string
	}
	cases := map[string]struct {
		args args
		want map[string]any
	}{
		"InjectsExternalName": {
			args: args{base: map[string]any{ProjectID: "p"}, externalName: "mdb_sa_id_1"},
			want: map[string]any{ProjectID: "p", ClientID: "mdb_sa_id_1"},
		},
		"OverridesExistingValue": {
			args: args{base: map[string]any{ClientID: "old"}, externalName: "new"},
			want: map[string]any{ClientID: "new"},
		},
		"SkipsEmptyExternalName": {
			args: args{base: map[string]any{ProjectID: "p"}, externalName: ""},
			want: map[string]any{ProjectID: "p"},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			SetIdentifierArgument(ClientID)(tc.args.base, tc.args.externalName)
			if diff := cmp.Diff(tc.want, tc.args.base); diff != "" {
				t.Errorf("SetIdentifierArgument(...): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestStateEmptyWhenAttributeUnset(t *testing.T) {
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{ClientID: tftypes.String, "name": tftypes.String}}
	state := func(clientID tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{
			ClientID: clientID,
			"name":   tftypes.NewValue(tftypes.String, "sa"),
		})
	}
	cases := map[string]struct {
		state tftypes.Value
		want  bool
	}{
		"NullState":        {state: tftypes.NewValue(objType, nil), want: true},
		"NullAttribute":    {state: state(tftypes.NewValue(tftypes.String, nil)), want: true},
		"UnknownAttribute": {state: state(tftypes.NewValue(tftypes.String, tftypes.UnknownValue)), want: true},
		"EmptyAttribute":   {state: state(tftypes.NewValue(tftypes.String, "")), want: true},
		"SetAttribute":     {state: state(tftypes.NewValue(tftypes.String, "mdb_sa_id_1")), want: false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := StateEmptyWhenAttributeUnset(ClientID)(context.Background(), tc.state, rschema.Schema{})
			if err != nil {
				t.Fatalf("StateEmptyWhenAttributeUnset(...): unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StateEmptyWhenAttributeUnset(...): -want, +got:\n%s", diff)
			}
		})
	}
}

const (
	cloudProviderConfigCount = "cloud_provider_config.#"
	dataProcessRegionCount   = "data_process_region.#"
)

func TestIgnoreUnknownUnsetBlocks(t *testing.T) {
	unknown := &terraform.ResourceAttrDiff{NewComputed: true}
	changed := &terraform.ResourceAttrDiff{Old: "a", New: "b"}
	existing := &terraform.InstanceState{ID: "id"}
	type args struct {
		diff   *terraform.InstanceDiff
		state  *terraform.InstanceState
		config *terraform.ResourceConfig
	}
	cases := map[string]struct {
		args args
		want *terraform.InstanceDiff
	}{
		"DropsUnknownUnsetBlock": {
			args: args{
				diff:   &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{cloudProviderConfigCount: unknown, dataProcessRegionCount: unknown}},
				state:  existing,
				config: terraform.NewResourceConfigRaw(map[string]any{Name: "x"}),
			},
			want: &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{}},
		},
		"KeepsOtherChanges": {
			args: args{
				diff:   &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{cloudProviderConfigCount: unknown, Name: changed}},
				state:  existing,
				config: terraform.NewResourceConfigRaw(map[string]any{Name: "b"}),
			},
			want: &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{Name: changed}},
		},
		"KeepsConfiguredBlock": {
			args: args{
				diff:   &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{dataProcessRegionCount: unknown}},
				state:  existing,
				config: terraform.NewResourceConfigRaw(map[string]any{"data_process_region": []any{map[string]any{"cloud_provider": "AWS"}}}),
			},
			want: &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{dataProcessRegionCount: unknown}},
		},
		"KeepsBlockSetInState": {
			args: args{
				diff:  &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{cloudProviderConfigCount: {Old: "1", NewComputed: true}}},
				state: existing,
			},
			want: &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{cloudProviderConfigCount: {Old: "1", NewComputed: true}}},
		},
		"KeepsCreateDiff": {
			args: args{
				diff:  &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{cloudProviderConfigCount: unknown}},
				state: &terraform.InstanceState{},
			},
			want: &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{cloudProviderConfigCount: unknown}},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := IgnoreUnknownUnsetBlocks("cloud_provider_config", "data_process_region")(tc.args.diff, tc.args.state, tc.args.config)
			if err != nil {
				t.Fatalf("IgnoreUnknownUnsetBlocks(...): unexpected error: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("IgnoreUnknownUnsetBlocks(...): -want, +got:\n%s", diff)
			}
		})
	}
}
