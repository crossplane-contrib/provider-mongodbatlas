/*
Copyright 2021 The Crossplane Authors.

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

// Code generated by upjet. DO NOT EDIT.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type SettingsOrgConfigObservation struct {
	ID *string `json:"id,omitempty" tf:"id,omitempty"`
}

type SettingsOrgConfigParameters struct {

	// +kubebuilder:validation:Optional
	DomainAllowList []*string `json:"domainAllowList,omitempty" tf:"domain_allow_list,omitempty"`

	// +kubebuilder:validation:Required
	DomainRestrictionEnabled *bool `json:"domainRestrictionEnabled" tf:"domain_restriction_enabled,omitempty"`

	// +kubebuilder:validation:Required
	FederationSettingsID *string `json:"federationSettingsId" tf:"federation_settings_id,omitempty"`

	// +kubebuilder:validation:Required
	IdentityProviderID *string `json:"identityProviderId" tf:"identity_provider_id,omitempty"`

	// +kubebuilder:validation:Required
	OrgID *string `json:"orgId" tf:"org_id,omitempty"`

	// +kubebuilder:validation:Optional
	PostAuthRoleGrants []*string `json:"postAuthRoleGrants,omitempty" tf:"post_auth_role_grants,omitempty"`
}

// SettingsOrgConfigSpec defines the desired state of SettingsOrgConfig
type SettingsOrgConfigSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     SettingsOrgConfigParameters `json:"forProvider"`
}

// SettingsOrgConfigStatus defines the observed state of SettingsOrgConfig.
type SettingsOrgConfigStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        SettingsOrgConfigObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// SettingsOrgConfig is the Schema for the SettingsOrgConfigs API. <no value>
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,mongodbatlas}
type SettingsOrgConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SettingsOrgConfigSpec   `json:"spec"`
	Status            SettingsOrgConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SettingsOrgConfigList contains a list of SettingsOrgConfigs
type SettingsOrgConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SettingsOrgConfig `json:"items"`
}

// Repository type metadata.
var (
	SettingsOrgConfig_Kind             = "SettingsOrgConfig"
	SettingsOrgConfig_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: SettingsOrgConfig_Kind}.String()
	SettingsOrgConfig_KindAPIVersion   = SettingsOrgConfig_Kind + "." + CRDGroupVersion.String()
	SettingsOrgConfig_GroupVersionKind = CRDGroupVersion.WithKind(SettingsOrgConfig_Kind)
)

func init() {
	SchemeBuilder.Register(&SettingsOrgConfig{}, &SettingsOrgConfigList{})
}
