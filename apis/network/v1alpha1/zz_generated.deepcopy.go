//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/apis/common/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Container) DeepCopyInto(out *Container) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Container.
func (in *Container) DeepCopy() *Container {
	if in == nil {
		return nil
	}
	out := new(Container)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Container) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerList) DeepCopyInto(out *ContainerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerList.
func (in *ContainerList) DeepCopy() *ContainerList {
	if in == nil {
		return nil
	}
	out := new(ContainerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ContainerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerObservation) DeepCopyInto(out *ContainerObservation) {
	*out = *in
	if in.AzureSubscriptionID != nil {
		in, out := &in.AzureSubscriptionID, &out.AzureSubscriptionID
		*out = new(string)
		**out = **in
	}
	if in.ContainerID != nil {
		in, out := &in.ContainerID, &out.ContainerID
		*out = new(string)
		**out = **in
	}
	if in.GCPProjectID != nil {
		in, out := &in.GCPProjectID, &out.GCPProjectID
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.NetworkName != nil {
		in, out := &in.NetworkName, &out.NetworkName
		*out = new(string)
		**out = **in
	}
	if in.Provisioned != nil {
		in, out := &in.Provisioned, &out.Provisioned
		*out = new(bool)
		**out = **in
	}
	if in.VPCID != nil {
		in, out := &in.VPCID, &out.VPCID
		*out = new(string)
		**out = **in
	}
	if in.VnetName != nil {
		in, out := &in.VnetName, &out.VnetName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerObservation.
func (in *ContainerObservation) DeepCopy() *ContainerObservation {
	if in == nil {
		return nil
	}
	out := new(ContainerObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerParameters) DeepCopyInto(out *ContainerParameters) {
	*out = *in
	if in.AtlasCidrBlock != nil {
		in, out := &in.AtlasCidrBlock, &out.AtlasCidrBlock
		*out = new(string)
		**out = **in
	}
	if in.ProjectID != nil {
		in, out := &in.ProjectID, &out.ProjectID
		*out = new(string)
		**out = **in
	}
	if in.ProjectIDRef != nil {
		in, out := &in.ProjectIDRef, &out.ProjectIDRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.ProjectIDSelector != nil {
		in, out := &in.ProjectIDSelector, &out.ProjectIDSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.ProviderName != nil {
		in, out := &in.ProviderName, &out.ProviderName
		*out = new(string)
		**out = **in
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
	if in.RegionName != nil {
		in, out := &in.RegionName, &out.RegionName
		*out = new(string)
		**out = **in
	}
	if in.Regions != nil {
		in, out := &in.Regions, &out.Regions
		*out = make([]*string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(string)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerParameters.
func (in *ContainerParameters) DeepCopy() *ContainerParameters {
	if in == nil {
		return nil
	}
	out := new(ContainerParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerSpec) DeepCopyInto(out *ContainerSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerSpec.
func (in *ContainerSpec) DeepCopy() *ContainerSpec {
	if in == nil {
		return nil
	}
	out := new(ContainerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerStatus) DeepCopyInto(out *ContainerStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerStatus.
func (in *ContainerStatus) DeepCopy() *ContainerStatus {
	if in == nil {
		return nil
	}
	out := new(ContainerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Peering) DeepCopyInto(out *Peering) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Peering.
func (in *Peering) DeepCopy() *Peering {
	if in == nil {
		return nil
	}
	out := new(Peering)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Peering) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeeringList) DeepCopyInto(out *PeeringList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Peering, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeeringList.
func (in *PeeringList) DeepCopy() *PeeringList {
	if in == nil {
		return nil
	}
	out := new(PeeringList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PeeringList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeeringObservation) DeepCopyInto(out *PeeringObservation) {
	*out = *in
	if in.AtlasID != nil {
		in, out := &in.AtlasID, &out.AtlasID
		*out = new(string)
		**out = **in
	}
	if in.ConnectionID != nil {
		in, out := &in.ConnectionID, &out.ConnectionID
		*out = new(string)
		**out = **in
	}
	if in.ErrorMessage != nil {
		in, out := &in.ErrorMessage, &out.ErrorMessage
		*out = new(string)
		**out = **in
	}
	if in.ErrorState != nil {
		in, out := &in.ErrorState, &out.ErrorState
		*out = new(string)
		**out = **in
	}
	if in.ErrorStateName != nil {
		in, out := &in.ErrorStateName, &out.ErrorStateName
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.PeerID != nil {
		in, out := &in.PeerID, &out.PeerID
		*out = new(string)
		**out = **in
	}
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(string)
		**out = **in
	}
	if in.StatusName != nil {
		in, out := &in.StatusName, &out.StatusName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeeringObservation.
func (in *PeeringObservation) DeepCopy() *PeeringObservation {
	if in == nil {
		return nil
	}
	out := new(PeeringObservation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeeringParameters) DeepCopyInto(out *PeeringParameters) {
	*out = *in
	if in.AccepterRegionName != nil {
		in, out := &in.AccepterRegionName, &out.AccepterRegionName
		*out = new(string)
		**out = **in
	}
	if in.AtlasCidrBlock != nil {
		in, out := &in.AtlasCidrBlock, &out.AtlasCidrBlock
		*out = new(string)
		**out = **in
	}
	if in.AtlasGCPProjectID != nil {
		in, out := &in.AtlasGCPProjectID, &out.AtlasGCPProjectID
		*out = new(string)
		**out = **in
	}
	if in.AtlasVPCName != nil {
		in, out := &in.AtlasVPCName, &out.AtlasVPCName
		*out = new(string)
		**out = **in
	}
	if in.AwsAccountID != nil {
		in, out := &in.AwsAccountID, &out.AwsAccountID
		*out = new(string)
		**out = **in
	}
	if in.AzureDirectoryID != nil {
		in, out := &in.AzureDirectoryID, &out.AzureDirectoryID
		*out = new(string)
		**out = **in
	}
	if in.AzureSubscriptionID != nil {
		in, out := &in.AzureSubscriptionID, &out.AzureSubscriptionID
		*out = new(string)
		**out = **in
	}
	if in.ContainerID != nil {
		in, out := &in.ContainerID, &out.ContainerID
		*out = new(string)
		**out = **in
	}
	if in.GCPProjectID != nil {
		in, out := &in.GCPProjectID, &out.GCPProjectID
		*out = new(string)
		**out = **in
	}
	if in.NetworkName != nil {
		in, out := &in.NetworkName, &out.NetworkName
		*out = new(string)
		**out = **in
	}
	if in.ProjectID != nil {
		in, out := &in.ProjectID, &out.ProjectID
		*out = new(string)
		**out = **in
	}
	if in.ProjectIDRef != nil {
		in, out := &in.ProjectIDRef, &out.ProjectIDRef
		*out = new(v1.Reference)
		**out = **in
	}
	if in.ProjectIDSelector != nil {
		in, out := &in.ProjectIDSelector, &out.ProjectIDSelector
		*out = new(v1.Selector)
		(*in).DeepCopyInto(*out)
	}
	if in.ProviderName != nil {
		in, out := &in.ProviderName, &out.ProviderName
		*out = new(string)
		**out = **in
	}
	if in.ResourceGroupName != nil {
		in, out := &in.ResourceGroupName, &out.ResourceGroupName
		*out = new(string)
		**out = **in
	}
	if in.RouteTableCidrBlock != nil {
		in, out := &in.RouteTableCidrBlock, &out.RouteTableCidrBlock
		*out = new(string)
		**out = **in
	}
	if in.VPCID != nil {
		in, out := &in.VPCID, &out.VPCID
		*out = new(string)
		**out = **in
	}
	if in.VnetName != nil {
		in, out := &in.VnetName, &out.VnetName
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeeringParameters.
func (in *PeeringParameters) DeepCopy() *PeeringParameters {
	if in == nil {
		return nil
	}
	out := new(PeeringParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeeringSpec) DeepCopyInto(out *PeeringSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeeringSpec.
func (in *PeeringSpec) DeepCopy() *PeeringSpec {
	if in == nil {
		return nil
	}
	out := new(PeeringSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PeeringStatus) DeepCopyInto(out *PeeringStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.AtProvider.DeepCopyInto(&out.AtProvider)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PeeringStatus.
func (in *PeeringStatus) DeepCopy() *PeeringStatus {
	if in == nil {
		return nil
	}
	out := new(PeeringStatus)
	in.DeepCopyInto(out)
	return out
}
