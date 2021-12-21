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

type TeamObservation struct {
	ID *string `json:"id,omitempty" tf:"id,omitempty"`

	TeamID *string `json:"teamId,omitempty" tf:"team_id,omitempty"`
}

type TeamParameters struct {

	// +kubebuilder:validation:Required
	OrgID *string `json:"orgId" tf:"org_id,omitempty"`

	// +kubebuilder:validation:Required
	Usernames []*string `json:"usernames" tf:"usernames,omitempty"`
}

// TeamSpec defines the desired state of Team
type TeamSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     TeamParameters `json:"forProvider"`
}

// TeamStatus defines the observed state of Team.
type TeamStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        TeamObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// Team is the Schema for the Teams API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,mongodbatlasjet}
type Team struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TeamSpec   `json:"spec"`
	Status            TeamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TeamList contains a list of Teams
type TeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Team `json:"items"`
}

// Repository type metadata.
var (
	Team_Kind             = "Team"
	Team_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: Team_Kind}.String()
	Team_KindAPIVersion   = Team_Kind + "." + CRDGroupVersion.String()
	Team_GroupVersionKind = CRDGroupVersion.WithKind(Team_Kind)
)

func init() {
	SchemeBuilder.Register(&Team{}, &TeamList{})
}
