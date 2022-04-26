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

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/terrajet/pkg/controller"

	configuration "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/alert/configuration"
	backupschedule "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupsnapshot"
	backupsnapshotrestorejob "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupsnapshotrestorejob"
	provideraccess "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/provideraccess"
	provideraccessauthorization "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/provideraccesssetup"
	providersnapshot "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/providersnapshot"
	providersnapshotbackuppolicy "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/providersnapshotbackuppolicy"
	providersnapshotrestorejob "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/providersnapshotrestorejob"
	dbrole "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/custom/dbrole"
	dnsconfigurationclusteraws "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/custom/dnsconfigurationclusteraws"
	lake "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/data/lake"
	user "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/database/user"
	trigger "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/event/trigger"
	clusterconfig "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/ldap/verify"
	window "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/maintenance/window"
	advancedcluster "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/mongodbatlas/cluster"
	project "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/mongodbatlas/project"
	team "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/network/container"
	peering "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/network/peering"
	archive "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/online/archive"
	invitation "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/org/invitation"
	ipmode "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/private/ipmode"
	endpoint "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpoint"
	endpointservice "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpointservice"
	endpointserviceadl "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpointserviceadl"
	invitationproject "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/project/ipaccesslist"
	providerconfig "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/providerconfig"
	index "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/search/index"
	partyintegration "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/third/partyintegration"
	authenticationdatabaseuser "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/x509/authenticationdatabaseuser"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		configuration.Setup,
		backupschedule.Setup,
		backupsnapshot.Setup,
		backupsnapshotrestorejob.Setup,
		provideraccess.Setup,
		provideraccessauthorization.Setup,
		provideraccesssetup.Setup,
		providersnapshot.Setup,
		providersnapshotbackuppolicy.Setup,
		providersnapshotrestorejob.Setup,
		dbrole.Setup,
		dnsconfigurationclusteraws.Setup,
		lake.Setup,
		user.Setup,
		trigger.Setup,
		clusterconfig.Setup,
		configurationldap.Setup,
		verify.Setup,
		window.Setup,
		advancedcluster.Setup,
		auditing.Setup,
		cluster.Setup,
		project.Setup,
		team.Setup,
		container.Setup,
		peering.Setup,
		archive.Setup,
		invitation.Setup,
		ipmode.Setup,
		endpoint.Setup,
		endpointservice.Setup,
		endpointserviceadl.Setup,
		invitationproject.Setup,
		ipaccesslist.Setup,
		providerconfig.Setup,
		index.Setup,
		partyintegration.Setup,
		authenticationdatabaseuser.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
