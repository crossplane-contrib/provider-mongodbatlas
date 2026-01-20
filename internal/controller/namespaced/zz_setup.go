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
	keyprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/api/keyprojectassignment"
	compliancepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/backup/compliancepolicy"
	backupschedule "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshot"
	backupsnapshotexportbucket "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotexportbucket"
	backupsnapshotexportjob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotexportjob"
	backupsnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotrestorejob"
	provideraccessauthorization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccesssetup"
	userorgassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/userorgassignment"
	userprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/userprojectassignment"
	userteamassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/userteamassignment"
	outagesimulation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cluster/outagesimulation"
	dbrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/custom/dbrole"
	dnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/custom/dnsconfigurationclusteraws"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/database/user"
	trigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/event/trigger"
	databaseinstance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/databaseinstance"
	querylimit "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/querylimit"
	settingsidentityprovider "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/settingsidentityprovider"
	settingsorgconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/settingsorgconfig"
	settingsorgrolemapping "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federated/settingsorgrolemapping"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/flex/cluster"
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/verify"
	window "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/maintenance/window"
	employeeaccessgrant "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodb/employeeaccessgrant"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/auditing"
	clustermongodbatlas "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/cluster"
	organization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/organization"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/project"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/peering"
	archive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/online/archive"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/org/invitation"
	endpointregionalmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/private/endpointregionalmode"
	endpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpoint"
	endpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointservice"
	endpointservicedatafederationonlinearchive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointservicedatafederationonlinearchive"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/ipaccesslist"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/providerconfig"
	basedlogexport "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/push/basedlogexport"
	policy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/resource/policy"
	deployment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/search/deployment"
	index "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/search/index"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/serverless/instance"
	connection "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/connection"
	instancestream "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/instance"
	privatelinkendpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/privatelinkendpoint"
	processor "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/processor"
	workspace "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/workspace"
	projectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/team/projectassignment"
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
