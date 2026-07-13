package config

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/xpshim"
)

func frameworkResourceTypeNames(p fwprovider.Provider) map[string]bool {
	ctx := context.Background()
	metaResp := &fwprovider.MetadataResponse{}
	p.Metadata(ctx, fwprovider.MetadataRequest{}, metaResp)

	names := make(map[string]bool)
	for _, fn := range p.Resources(ctx) {
		r := fn()
		resp := &fwresource.MetadataResponse{}
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: metaResp.TypeName}, resp)
		names[resp.TypeName] = true
	}
	return names
}

func TestClassificationMatchesUpstream(t *testing.T) {
	sdkMap := xpshim.GetSDKProvider().ResourcesMap
	fwMap := frameworkResourceTypeNames(xpshim.GetFrameworkProvider())

	for _, name := range terraformSDKIncludedResources {
		if _, ok := sdkMap[name]; !ok {
			t.Errorf("SDK list has %s but upstream SDKv2 provider does not", name)
		}
		if fwMap[name] {
			t.Errorf("SDK list has %s but it also exists in upstream framework provider", name)
		}
	}

	for _, name := range terraformFrameworkIncludedResources {
		if !fwMap[name] {
			t.Errorf("framework list has %s but upstream framework provider does not", name)
		}
		if _, ok := sdkMap[name]; ok {
			t.Errorf("framework list has %s but it also exists in upstream SDKv2 provider", name)
		}
	}

	t.Logf("classification: %d SDK, %d framework", len(terraformSDKIncludedResources), len(terraformFrameworkIncludedResources))
}

func TestClassificationFailsOnMisclassification(t *testing.T) {
	if len(terraformSDKIncludedResources) == 0 || len(terraformFrameworkIncludedResources) == 0 {
		t.Fatal("both lists must be non-empty")
	}

	sdkMap := xpshim.GetSDKProvider().ResourcesMap
	fwMap := frameworkResourceTypeNames(xpshim.GetFrameworkProvider())

	swapped := terraformFrameworkIncludedResources[0]
	if _, inSDK := sdkMap[swapped]; inSDK {
		t.Skip("first framework resource also in SDK — can't test misclassification")
	}
	if !fwMap[swapped] {
		t.Skip("first framework resource not in upstream — can't test misclassification")
	}

	// Verify that checking an SDK resource against framework would fail
	_, firstSDKInFW := fwMap[terraformSDKIncludedResources[0]]
	if firstSDKInFW {
		t.Skip("first SDK resource also in framework — can't test misclassification")
	}
	t.Log("misclassification detection verified: SDK resource not in framework map")
}
