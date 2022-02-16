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

type EndpointServiceObservation struct {
	AwsConnectionStatus *string `json:"awsConnectionStatus,omitempty" tf:"aws_connection_status,omitempty"`

	AzureStatus *string `json:"azureStatus,omitempty" tf:"azure_status,omitempty"`

	DeleteRequested *bool `json:"deleteRequested,omitempty" tf:"delete_requested,omitempty"`

	EndpointGroupName *string `json:"endpointGroupName,omitempty" tf:"endpoint_group_name,omitempty"`

	ErrorMessage *string `json:"errorMessage,omitempty" tf:"error_message,omitempty"`

	GCPStatus *string `json:"gcpStatus,omitempty" tf:"gcp_status,omitempty"`

	ID *string `json:"id,omitempty" tf:"id,omitempty"`

	InterfaceEndpointID *string `json:"interfaceEndpointId,omitempty" tf:"interface_endpoint_id,omitempty"`

	PrivateEndpointConnectionName *string `json:"privateEndpointConnectionName,omitempty" tf:"private_endpoint_connection_name,omitempty"`

	PrivateEndpointResourceID *string `json:"privateEndpointResourceId,omitempty" tf:"private_endpoint_resource_id,omitempty"`
}

type EndpointServiceParameters struct {

	// +kubebuilder:validation:Required
	EndpointServiceID *string `json:"endpointServiceId" tf:"endpoint_service_id,omitempty"`

	// +kubebuilder:validation:Optional
	Endpoints []EndpointsParameters `json:"endpoints,omitempty" tf:"endpoints,omitempty"`

	// +kubebuilder:validation:Optional
	GCPProjectID *string `json:"gcpProjectId,omitempty" tf:"gcp_project_id,omitempty"`

	// +kubebuilder:validation:Optional
	PrivateEndpointIPAddress *string `json:"privateEndpointIpAddress,omitempty" tf:"private_endpoint_ip_address,omitempty"`

	// +kubebuilder:validation:Required
	PrivateLinkID *string `json:"privateLinkId" tf:"private_link_id,omitempty"`

	// +crossplane:generate:reference:type=github.com/crossplane-contrib/provider-jet-mongodbatlas/apis/mongodbatlas/v1alpha1.Project
	// +crossplane:generate:reference:extractor=github.com/crossplane-contrib/provider-jet-mongodbatlas/config/common.ExtractResourceID()
	// +kubebuilder:validation:Optional
	ProjectID *string `json:"projectId,omitempty" tf:"project_id,omitempty"`

	// +kubebuilder:validation:Optional
	ProjectIDRef *v1.Reference `json:"projectIdRef,omitempty" tf:"-"`

	// +kubebuilder:validation:Optional
	ProjectIDSelector *v1.Selector `json:"projectIdSelector,omitempty" tf:"-"`

	// +kubebuilder:validation:Required
	ProviderName *string `json:"providerName" tf:"provider_name,omitempty"`
}

type EndpointsObservation struct {
	ServiceAttachmentName *string `json:"serviceAttachmentName,omitempty" tf:"service_attachment_name,omitempty"`

	Status *string `json:"status,omitempty" tf:"status,omitempty"`
}

type EndpointsParameters struct {

	// +kubebuilder:validation:Optional
	EndpointName *string `json:"endpointName,omitempty" tf:"endpoint_name,omitempty"`

	// +kubebuilder:validation:Optional
	IPAddress *string `json:"ipAddress,omitempty" tf:"ip_address,omitempty"`
}

// EndpointServiceSpec defines the desired state of EndpointService
type EndpointServiceSpec struct {
	v1.ResourceSpec `json:",inline"`
	ForProvider     EndpointServiceParameters `json:"forProvider"`
}

// EndpointServiceStatus defines the observed state of EndpointService.
type EndpointServiceStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        EndpointServiceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// EndpointService is the Schema for the EndpointServices API
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,mongodbatlasjet}
type EndpointService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              EndpointServiceSpec   `json:"spec"`
	Status            EndpointServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EndpointServiceList contains a list of EndpointServices
type EndpointServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EndpointService `json:"items"`
}

// Repository type metadata.
var (
	EndpointService_Kind             = "EndpointService"
	EndpointService_GroupKind        = schema.GroupKind{Group: CRDGroup, Kind: EndpointService_Kind}.String()
	EndpointService_KindAPIVersion   = EndpointService_Kind + "." + CRDGroupVersion.String()
	EndpointService_GroupVersionKind = CRDGroupVersion.WithKind(EndpointService_Kind)
)

func init() {
	SchemeBuilder.Register(&EndpointService{}, &EndpointServiceList{})
}
