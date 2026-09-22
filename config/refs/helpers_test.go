package refs

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

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
