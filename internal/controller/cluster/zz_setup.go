// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	listapikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/access/listapikey"
	configuration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/alert/configuration"
	key "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/api/key"
	compliancepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/backup/compliancepolicy"
	backupschedule "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshot"
	backupsnapshotexportbucket "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshotexportbucket"
	backupsnapshotexportjob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshotexportjob"
	backupsnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshotrestorejob"
	provideraccess "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccess"
	provideraccessauthorization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccesssetup"
	providersnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/providersnapshot"
	providersnapshotbackuppolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/providersnapshotbackuppolicy"
	providersnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/providersnapshotrestorejob"
	dbrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/custom/dbrole"
	dnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/custom/dnsconfigurationclusteraws"
	lake "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/data/lake"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/user"
	trigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/event/trigger"
	settingsidentityprovider "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsidentityprovider"
	settingsorgconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsorgconfig"
	settingsorgrolemapping "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsorgrolemapping"
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/verify"
	window "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/maintenance/window"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/cluster"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/project"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/peering"
	archive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/online/archive"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/org/invitation"
	endpointregionalmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/private/endpointregionalmode"
	ipmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/private/ipmode"
	endpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpoint"
	endpointserverless "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpointserverless"
	endpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpointservice"
	endpointserviceadl "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpointserviceadl"
	endpointserviceserverless "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpointserviceserverless"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/ipaccesslist"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/providerconfig"
	index "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/search/index"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/serverless/instance"
	partyintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/third/partyintegration"
	authenticationdatabaseuser "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/x509/authenticationdatabaseuser"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		listapikey.Setup,
		configuration.Setup,
		key.Setup,
		compliancepolicy.Setup,
		backupschedule.Setup,
		backupsnapshot.Setup,
		backupsnapshotexportbucket.Setup,
		backupsnapshotexportjob.Setup,
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
		settingsidentityprovider.Setup,
		settingsorgconfig.Setup,
		settingsorgrolemapping.Setup,
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
		endpointregionalmode.Setup,
		ipmode.Setup,
		endpoint.Setup,
		endpointserverless.Setup,
		endpointservice.Setup,
		endpointserviceadl.Setup,
		endpointserviceserverless.Setup,
		apikey.Setup,
		invitationproject.Setup,
		ipaccesslist.Setup,
		providerconfig.Setup,
		index.Setup,
		instance.Setup,
		partyintegration.Setup,
		authenticationdatabaseuser.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		listapikey.SetupGated,
		configuration.SetupGated,
		key.SetupGated,
		compliancepolicy.SetupGated,
		backupschedule.SetupGated,
		backupsnapshot.SetupGated,
		backupsnapshotexportbucket.SetupGated,
		backupsnapshotexportjob.SetupGated,
		backupsnapshotrestorejob.SetupGated,
		provideraccess.SetupGated,
		provideraccessauthorization.SetupGated,
		provideraccesssetup.SetupGated,
		providersnapshot.SetupGated,
		providersnapshotbackuppolicy.SetupGated,
		providersnapshotrestorejob.SetupGated,
		dbrole.SetupGated,
		dnsconfigurationclusteraws.SetupGated,
		lake.SetupGated,
		user.SetupGated,
		trigger.SetupGated,
		settingsidentityprovider.SetupGated,
		settingsorgconfig.SetupGated,
		settingsorgrolemapping.SetupGated,
		clusterconfig.SetupGated,
		configurationldap.SetupGated,
		verify.SetupGated,
		window.SetupGated,
		advancedcluster.SetupGated,
		auditing.SetupGated,
		cluster.SetupGated,
		project.SetupGated,
		team.SetupGated,
		container.SetupGated,
		peering.SetupGated,
		archive.SetupGated,
		invitation.SetupGated,
		endpointregionalmode.SetupGated,
		ipmode.SetupGated,
		endpoint.SetupGated,
		endpointserverless.SetupGated,
		endpointservice.SetupGated,
		endpointserviceadl.SetupGated,
		endpointserviceserverless.SetupGated,
		apikey.SetupGated,
		invitationproject.SetupGated,
		ipaccesslist.SetupGated,
		providerconfig.SetupGated,
		index.SetupGated,
		instance.SetupGated,
		partyintegration.SetupGated,
		authenticationdatabaseuser.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
