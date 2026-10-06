package refs

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
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
