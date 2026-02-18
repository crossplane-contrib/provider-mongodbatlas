// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	listapikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/access/listapikey"
	configuration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/alert/configuration"
	compliancepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/backup/compliancepolicy"
	backupschedule "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshot"
	backupsnapshotexportbucket "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotexportbucket"
	backupsnapshotexportjob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotexportjob"
	backupsnapshotrestorejob "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/backupsnapshotrestorejob"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/instance"
	provideraccessauthorization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/provideraccesssetup"
	userorgassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/userorgassignment"
	userprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/userprojectassignment"
	userteamassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/cloud/userteamassignment"
	customrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/database/customrole"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/database/user"
	x509userauthentication "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/database/x509userauthentication"
	databaseinstance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federation/databaseinstance"
	identityprovider "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federation/identityprovider"
	orgconfigsettings "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federation/orgconfigsettings"
	privatelinkendpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federation/privatelinkendpointservice"
	querylimit "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federation/querylimit"
	rolemapping "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/federation/rolemapping"
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/verify"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/advancedcluster"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/apikey"
	apikeyprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/apikeyprojectassignment"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/cluster"
	clusteroutagesimulation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/clusteroutagesimulation"
	customdnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/customdnsconfigurationclusteraws"
	employeeaccessgrant "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/employeeaccessgrant"
	eventtrigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/eventtrigger"
	flexcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/flexcluster"
	logintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/logintegration"
	maintenancewindow "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/maintenancewindow"
	onlinearchive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/onlinearchive"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/project"
	pushbasedlogexport "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/pushbasedlogexport"
	resourcepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/resourcepolicy"
	serviceaccount "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/serviceaccount"
	serviceaccountaccesslistentry "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/serviceaccountaccesslistentry"
	serviceaccountprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/serviceaccountprojectassignment"
	serviceaccountsecret "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/serviceaccountsecret"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/peering"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/org/invitation"
	organization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/org/organization"
	regionalmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privateendpoint/regionalmode"
	resource "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privateendpoint/resource"
	service "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privateendpoint/service"
	apikeyproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/ipaccesslist"
	serviceaccountproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/serviceaccount"
	serviceaccountaccesslistentryproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/serviceaccountaccesslistentry"
	serviceaccountsecretproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/serviceaccountsecret"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/providerconfig"
	deployment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/search/deployment"
	index "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/search/index"
	connection "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/connection"
	instancestream "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/instance"
	privatelinkendpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/privatelinkendpoint"
	processor "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/processor"
	workspace "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/stream/workspace"
	projectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/team/projectassignment"
	partyintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/third/partyintegration"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		listapikey.Setup,
		configuration.Setup,
		compliancepolicy.Setup,
		backupschedule.Setup,
		backupsnapshot.Setup,
		backupsnapshotexportbucket.Setup,
		backupsnapshotexportjob.Setup,
		backupsnapshotrestorejob.Setup,
		instance.Setup,
		provideraccessauthorization.Setup,
		provideraccesssetup.Setup,
		userorgassignment.Setup,
		userprojectassignment.Setup,
		userteamassignment.Setup,
		customrole.Setup,
		user.Setup,
		x509userauthentication.Setup,
		databaseinstance.Setup,
		identityprovider.Setup,
		orgconfigsettings.Setup,
		privatelinkendpointservice.Setup,
		querylimit.Setup,
		rolemapping.Setup,
		clusterconfig.Setup,
		configurationldap.Setup,
		verify.Setup,
		advancedcluster.Setup,
		apikey.Setup,
		apikeyprojectassignment.Setup,
		auditing.Setup,
		cluster.Setup,
		clusteroutagesimulation.Setup,
		customdnsconfigurationclusteraws.Setup,
		employeeaccessgrant.Setup,
		eventtrigger.Setup,
		flexcluster.Setup,
		logintegration.Setup,
		maintenancewindow.Setup,
		onlinearchive.Setup,
		project.Setup,
		pushbasedlogexport.Setup,
		resourcepolicy.Setup,
		serviceaccount.Setup,
		serviceaccountaccesslistentry.Setup,
		serviceaccountprojectassignment.Setup,
		serviceaccountsecret.Setup,
		team.Setup,
		container.Setup,
		peering.Setup,
		invitation.Setup,
		organization.Setup,
		regionalmode.Setup,
		resource.Setup,
		service.Setup,
		apikeyproject.Setup,
		invitationproject.Setup,
		ipaccesslist.Setup,
		serviceaccountproject.Setup,
		serviceaccountaccesslistentryproject.Setup,
		serviceaccountsecretproject.Setup,
		providerconfig.Setup,
		deployment.Setup,
		index.Setup,
		connection.Setup,
		instancestream.Setup,
		privatelinkendpoint.Setup,
		processor.Setup,
		workspace.Setup,
		projectassignment.Setup,
		partyintegration.Setup,
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
		compliancepolicy.SetupGated,
		backupschedule.SetupGated,
		backupsnapshot.SetupGated,
		backupsnapshotexportbucket.SetupGated,
		backupsnapshotexportjob.SetupGated,
		backupsnapshotrestorejob.SetupGated,
		instance.SetupGated,
		provideraccessauthorization.SetupGated,
		provideraccesssetup.SetupGated,
		userorgassignment.SetupGated,
		userprojectassignment.SetupGated,
		userteamassignment.SetupGated,
		customrole.SetupGated,
		user.SetupGated,
		x509userauthentication.SetupGated,
		databaseinstance.SetupGated,
		identityprovider.SetupGated,
		orgconfigsettings.SetupGated,
		privatelinkendpointservice.SetupGated,
		querylimit.SetupGated,
		rolemapping.SetupGated,
		clusterconfig.SetupGated,
		configurationldap.SetupGated,
		verify.SetupGated,
		advancedcluster.SetupGated,
		apikey.SetupGated,
		apikeyprojectassignment.SetupGated,
		auditing.SetupGated,
		cluster.SetupGated,
		clusteroutagesimulation.SetupGated,
		customdnsconfigurationclusteraws.SetupGated,
		employeeaccessgrant.SetupGated,
		eventtrigger.SetupGated,
		flexcluster.SetupGated,
		logintegration.SetupGated,
		maintenancewindow.SetupGated,
		onlinearchive.SetupGated,
		project.SetupGated,
		pushbasedlogexport.SetupGated,
		resourcepolicy.SetupGated,
		serviceaccount.SetupGated,
		serviceaccountaccesslistentry.SetupGated,
		serviceaccountprojectassignment.SetupGated,
		serviceaccountsecret.SetupGated,
		team.SetupGated,
		container.SetupGated,
		peering.SetupGated,
		invitation.SetupGated,
		organization.SetupGated,
		regionalmode.SetupGated,
		resource.SetupGated,
		service.SetupGated,
		apikeyproject.SetupGated,
		invitationproject.SetupGated,
		ipaccesslist.SetupGated,
		serviceaccountproject.SetupGated,
		serviceaccountaccesslistentryproject.SetupGated,
		serviceaccountsecretproject.SetupGated,
		providerconfig.SetupGated,
		deployment.SetupGated,
		index.SetupGated,
		connection.SetupGated,
		instancestream.SetupGated,
		privatelinkendpoint.SetupGated,
		processor.SetupGated,
		workspace.SetupGated,
		projectassignment.SetupGated,
		partyintegration.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
