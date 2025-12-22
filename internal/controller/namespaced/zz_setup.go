// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	listapikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/access/listapikey"
	configuration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/alert/configuration"
	key "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/api/key"
	compliancepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/backup/compliancepolicy"
	backupschedule "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshot"
	backupsnapshotexportbucket "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotexportbucket"
	backupsnapshotexportjob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotexportjob"
	backupsnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotrestorejob"
	provideraccess "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccess"
	provideraccessauthorization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccesssetup"
	providersnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/providersnapshot"
	providersnapshotbackuppolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/providersnapshotbackuppolicy"
	providersnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/providersnapshotrestorejob"
	dbrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/custom/dbrole"
	dnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/custom/dnsconfigurationclusteraws"
	lake "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/data/lake"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/database/user"
	trigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/event/trigger"
	settingsidentityprovider "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/settingsidentityprovider"
	settingsorgconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/settingsorgconfig"
	settingsorgrolemapping "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/settingsorgrolemapping"
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/verify"
	window "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/maintenance/window"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/cluster"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/project"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/peering"
	archive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/online/archive"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/org/invitation"
	endpointregionalmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/private/endpointregionalmode"
	ipmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/private/ipmode"
	endpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpoint"
	endpointserverless "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointserverless"
	endpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointservice"
	endpointserviceadl "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointserviceadl"
	endpointserviceserverless "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointserviceserverless"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/ipaccesslist"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/providerconfig"
	index "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/search/index"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/serverless/instance"
	partyintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/third/partyintegration"
	authenticationdatabaseuser "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/x509/authenticationdatabaseuser"
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
