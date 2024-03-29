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

type EndpointRegionalModeObservation struct {
	ID *string `json:"id,omitempty" tf:"id,omitempty"`
}

type EndpointRegionalModeParameters struct {

	// +kubebuilder:validation:Optional
	Enabled *bool `json:"enabled,omitempty" tf:"enabled,omitempty"`

	// +crossplane:generate:reference:type=github.com/crossplane-contrib/provider-mongodbatlas/apis/mongodbatlas/v1alpha1.Project
	// +crossplane:generate:reference:extractor=github.com/crossplane-contrib/provider-mongodbatlas/config/common.ExtractResourceID()
	// +kubebuilder:validation:Optional
	ProjectID *string `json:"projectId,omitempty" tf:"project_id,omitempty"`

	// Reference to a Project in mongodbatlas to populate projectId.
	// +kubebuilder:validation:Optional
	ProjectIDRef *v1.Reference `json:"projectIdRef,omitempty" tf:"-"`

	// Selector for a Project in mongodbatlas to populate projectId.
	// +kubebuilder:validation:Optional
	ProjectIDSelector *v1.Selector `json:"projectIdSelector,omitempty" tf:"-"`
}

// EndpointRegionalModeSpec defines the desired state of EndpointRegionalMode
type EndpointRegionalModeSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     EndpointRegionalModeParameters `json:"forProvider"`
}

// EndpointRegionalModeStatus defines the observed state of EndpointRegionalMode.
type EndpointRegionalModeStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        EndpointRegionalModeObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// EndpointRegionalMode is the Schema for the EndpointRegionalModes API. <no value>
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,mongodbatlas}
type EndpointRegionalMode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EndpointRegionalModeSpec   `json:"spec"`
	Status            EndpointRegionalModeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EndpointRegionalModeList contains a list of EndpointRegionalModes
type EndpointRegionalModeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EndpointRegionalMode `json:"items"`
}

// Repository type metadata.
var (
	EndpointRegionalMode_Kind             = "EndpointRegionalMode"
	EndpointRegionalMode_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: EndpointRegionalMode_Kind}.String()
	EndpointRegionalMode_KindAPIVersion   = EndpointRegionalMode_Kind + "." + CRDGroupVersion.String()
	EndpointRegionalMode_GroupVersionKind = CRDGroupVersion.WithKind(EndpointRegionalMode_Kind)
)

func init() {
	SchemeBuilder.Register(&EndpointRegionalMode{}, &EndpointRegionalModeList{})
}
