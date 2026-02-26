// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	configuration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/alert/configuration"
	backupcompliancepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/backupcompliancepolicy"
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
	customrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/customrole"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/user"
	x509userauthentication "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/x509userauthentication"
	databaseinstance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/databaseinstance"
	orgconfigsettings "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/orgconfigsettings"
	privatelinkendpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/privatelinkendpointservice"
	querylimit "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/querylimit"
	rolemapping "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/rolemapping"
	settingsidentityprovider "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federated/settingsidentityprovider"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/configuration"
	verify "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/verify"
	accesslistapikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/accesslistapikey"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/advancedcluster"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/apikey"
	apikeyprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/apikeyprojectassignment"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/cluster"
	customdnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/customdnsconfigurationclusteraws"
	employeeaccessgrant "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/employeeaccessgrant"
	eventtrigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/eventtrigger"
	flexcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/flexcluster"
	globalclusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/globalclusterconfig"
	logintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/logintegration"
	maintenancewindow "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/maintenancewindow"
	onlinearchive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/onlinearchive"
	organization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/organization"
	partyintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/partyintegration"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/project"
	pushbasedlogexport "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/pushbasedlogexport"
	resourcepolicy "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/resourcepolicy"
	serviceaccount "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/serviceaccount"
	serviceaccountaccesslistentry "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/serviceaccountaccesslistentry"
	serviceaccountprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/serviceaccountprojectassignment"
	serviceaccountsecret "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/serviceaccountsecret"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/peering"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/org/invitation"
	regionalmode "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privateendpoint/regionalmode"
	resource "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privateendpoint/resource"
	service "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/privateendpoint/service"
	apikeyproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/ipaccesslist"
	serviceaccountproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/serviceaccount"
	serviceaccountaccesslistentryproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/serviceaccountaccesslistentry"
	serviceaccountsecretproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/serviceaccountsecret"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/providerconfig"
	deployment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/search/deployment"
	index "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/search/index"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/serverless/instance"
	connection "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/connection"
	instancestream "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/instance"
	privatelinkendpoint "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/privatelinkendpoint"
	processor "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/processor"
	workspace "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/stream/workspace"
	projectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/team/projectassignment"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		configuration.Setup,
		backupcompliancepolicy.Setup,
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
		customrole.Setup,
		user.Setup,
		x509userauthentication.Setup,
		databaseinstance.Setup,
		orgconfigsettings.Setup,
		privatelinkendpointservice.Setup,
		querylimit.Setup,
		rolemapping.Setup,
		settingsidentityprovider.Setup,
		configurationldap.Setup,
		verify.Setup,
		accesslistapikey.Setup,
		advancedcluster.Setup,
		apikey.Setup,
		apikeyprojectassignment.Setup,
		auditing.Setup,
		cluster.Setup,
		customdnsconfigurationclusteraws.Setup,
		employeeaccessgrant.Setup,
		eventtrigger.Setup,
		flexcluster.Setup,
		globalclusterconfig.Setup,
		logintegration.Setup,
		maintenancewindow.Setup,
		onlinearchive.Setup,
		organization.Setup,
		partyintegration.Setup,
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
		instance.Setup,
		connection.Setup,
		instancestream.Setup,
		privatelinkendpoint.Setup,
		processor.Setup,
		workspace.Setup,
		projectassignment.Setup,
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
		configuration.SetupGated,
		backupcompliancepolicy.SetupGated,
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
		customrole.SetupGated,
		user.SetupGated,
		x509userauthentication.SetupGated,
		databaseinstance.SetupGated,
		orgconfigsettings.SetupGated,
		privatelinkendpointservice.SetupGated,
		querylimit.SetupGated,
		rolemapping.SetupGated,
		settingsidentityprovider.SetupGated,
		configurationldap.SetupGated,
		verify.SetupGated,
		accesslistapikey.SetupGated,
		advancedcluster.SetupGated,
		apikey.SetupGated,
		apikeyprojectassignment.SetupGated,
		auditing.SetupGated,
		cluster.SetupGated,
		customdnsconfigurationclusteraws.SetupGated,
		employeeaccessgrant.SetupGated,
		eventtrigger.SetupGated,
		flexcluster.SetupGated,
		globalclusterconfig.SetupGated,
		logintegration.SetupGated,
		maintenancewindow.SetupGated,
		onlinearchive.SetupGated,
		organization.SetupGated,
		partyintegration.SetupGated,
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
		instance.SetupGated,
		connection.SetupGated,
		instancestream.SetupGated,
		privatelinkendpoint.SetupGated,
		processor.SetupGated,
		workspace.SetupGated,
		projectassignment.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
