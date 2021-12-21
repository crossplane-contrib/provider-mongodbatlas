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

// Code generated by terrajet. DO NOT EDIT.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

type IPModeObservation struct {
	ID *string `json:"id,omitempty" tf:"id,omitempty"`
}

type IPModeParameters struct {

	// +kubebuilder:validation:Required
	Enabled *bool `json:"enabled" tf:"enabled,omitempty"`

	// +kubebuilder:validation:Required
	ProjectID *string `json:"projectId" tf:"project_id,omitempty"`
}

// IPModeSpec defines the desired state of IPMode
type IPModeSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     IPModeParameters `json:"forProvider"`
}

// IPModeStatus defines the observed state of IPMode.
type IPModeStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        IPModeObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// IPMode is the Schema for the IPModes API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,mongodbatlasjet}
type IPMode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IPModeSpec   `json:"spec"`
	Status            IPModeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IPModeList contains a list of IPModes
type IPModeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IPMode `json:"items"`
}

// Repository type metadata.
var (
	IPMode_Kind             = "IPMode"
	IPMode_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: IPMode_Kind}.String()
	IPMode_KindAPIVersion   = IPMode_Kind + "." + CRDGroupVersion.String()
	IPMode_GroupVersionKind = CRDGroupVersion.WithKind(IPMode_Kind)
)

func init() {
	SchemeBuilder.Register(&IPMode{}, &IPModeList{})
}
