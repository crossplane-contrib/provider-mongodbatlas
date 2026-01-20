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
	keyprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/api/keyprojectassignment"
	compliancepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/backup/compliancepolicy"
	backupschedule "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshot"
	backupsnapshotexportbucket "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshotexportbucket"
	backupsnapshotexportjob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshotexportjob"
	backupsnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupsnapshotrestorejob"
	provideraccessauthorization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccesssetup"
	userorgassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/userorgassignment"
	userprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/userprojectassignment"
	userteamassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/userteamassignment"
	outagesimulation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cluster/outagesimulation"
	dbrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/custom/dbrole"
	dnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/custom/dnsconfigurationclusteraws"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/user"
	trigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/event/trigger"
	databaseinstance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/databaseinstance"
	querylimit "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/querylimit"
	settingsidentityprovider "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsidentityprovider"
	settingsorgconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsorgconfig"
	settingsorgrolemapping "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsorgrolemapping"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/flex/cluster"
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/verify"
	window "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/maintenance/window"
	employeeaccessgrant "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodb/employeeaccessgrant"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/auditing"
	clustermongodbatlas "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/cluster"
	organization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/organization"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/project"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/peering"
	archive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/online/archive"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/org/invitation"
	endpointregionalmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/private/endpointregionalmode"
	endpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpoint"
	endpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpointservice"
	endpointservicedatafederationonlinearchive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privatelink/endpointservicedatafederationonlinearchive"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/ipaccesslist"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/providerconfig"
	basedlogexport "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/push/basedlogexport"
	policy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/resource/policy"
	deployment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/search/deployment"
	index "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/search/index"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/serverless/instance"
	connection "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/connection"
	instancestream "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/instance"
	privatelinkendpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/privatelinkendpoint"
	processor "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/processor"
	workspace "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/workspace"
	projectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/team/projectassignment"
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
		keyprojectassignment.Setup,
		compliancepolicy.Setup,
		backupschedule.Setup,
		backupsnapshot.Setup,
		backupsnapshotexportbucket.Setup,
		backupsnapshotexportjob.Setup,
		backupsnapshotrestorejob.Setup,
		provideraccessauthorization.Setup,
		provideraccesssetup.Setup,
		userorgassignment.Setup,
		userprojectassignment.Setup,
		userteamassignment.Setup,
		outagesimulation.Setup,
		dbrole.Setup,
		dnsconfigurationclusteraws.Setup,
		user.Setup,
		trigger.Setup,
		databaseinstance.Setup,
		querylimit.Setup,
		settingsidentityprovider.Setup,
		settingsorgconfig.Setup,
		settingsorgrolemapping.Setup,
		cluster.Setup,
		clusterconfig.Setup,
		configurationldap.Setup,
		verify.Setup,
		window.Setup,
		employeeaccessgrant.Setup,
		advancedcluster.Setup,
		auditing.Setup,
		clustermongodbatlas.Setup,
		organization.Setup,
		project.Setup,
		team.Setup,
		container.Setup,
		peering.Setup,
		archive.Setup,
		invitation.Setup,
		endpointregionalmode.Setup,
		endpoint.Setup,
		endpointservice.Setup,
		endpointservicedatafederationonlinearchive.Setup,
		apikey.Setup,
		invitationproject.Setup,
		ipaccesslist.Setup,
		providerconfig.Setup,
		basedlogexport.Setup,
		policy.Setup,
		deployment.Setup,
		index.Setup,
		instance.Setup,
		connection.Setup,
		instancestream.Setup,
		privatelinkendpoint.Setup,
		processor.Setup,
		workspace.Setup,
		projectassignment.Setup,
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
		keyprojectassignment.SetupGated,
		compliancepolicy.SetupGated,
		backupschedule.SetupGated,
		backupsnapshot.SetupGated,
		backupsnapshotexportbucket.SetupGated,
		backupsnapshotexportjob.SetupGated,
		backupsnapshotrestorejob.SetupGated,
		provideraccessauthorization.SetupGated,
		provideraccesssetup.SetupGated,
		userorgassignment.SetupGated,
		userprojectassignment.SetupGated,
		userteamassignment.SetupGated,
		outagesimulation.SetupGated,
		dbrole.SetupGated,
		dnsconfigurationclusteraws.SetupGated,
		user.SetupGated,
		trigger.SetupGated,
		databaseinstance.SetupGated,
		querylimit.SetupGated,
		settingsidentityprovider.SetupGated,
		settingsorgconfig.SetupGated,
		settingsorgrolemapping.SetupGated,
		cluster.SetupGated,
		clusterconfig.SetupGated,
		configurationldap.SetupGated,
		verify.SetupGated,
		window.SetupGated,
		employeeaccessgrant.SetupGated,
		advancedcluster.SetupGated,
		auditing.SetupGated,
		clustermongodbatlas.SetupGated,
		organization.SetupGated,
		project.SetupGated,
		team.SetupGated,
		container.SetupGated,
		peering.SetupGated,
		archive.SetupGated,
		invitation.SetupGated,
		endpointregionalmode.SetupGated,
		endpoint.SetupGated,
		endpointservice.SetupGated,
		endpointservicedatafederationonlinearchive.SetupGated,
		apikey.SetupGated,
		invitationproject.SetupGated,
		ipaccesslist.SetupGated,
		providerconfig.SetupGated,
		basedlogexport.SetupGated,
		policy.SetupGated,
		deployment.SetupGated,
		index.SetupGated,
		instance.SetupGated,
		connection.SetupGated,
		instancestream.SetupGated,
		privatelinkendpoint.SetupGated,
		processor.SetupGated,
		workspace.SetupGated,
		projectassignment.SetupGated,
		partyintegration.SetupGated,
		authenticationdatabaseuser.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
