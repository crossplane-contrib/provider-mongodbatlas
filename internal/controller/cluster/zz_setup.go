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
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/instance"
	provideraccessauthorization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccessauthorization"
	provideraccesssetup "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/provideraccesssetup"
	userorgassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/userorgassignment"
	userprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/userprojectassignment"
	userteamassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/cloud/userteamassignment"
	dnsconfigurationclusteraws "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/custom/dnsconfigurationclusteraws"
	customrole "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/customrole"
	user "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/user"
	x509userauthentication "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/database/x509userauthentication"
	trigger "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/event/trigger"
	databaseinstance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federation/databaseinstance"
	privatelinkendpointservice "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federation/privatelinkendpointservice"
	querylimit "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/federation/querylimit"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/ldap/configuration"
	integration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/log/integration"
	window "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/maintenance/window"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/cluster"
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/clusterconfig"
	clusteroutagesimulation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/clusteroutagesimulation"
	employeeaccessgrant "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/employeeaccessgrant"
	organization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/organization"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/project"
	serviceaccountaccesslistentry "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/serviceaccountaccesslistentry"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/network/peering"
	archive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/online/archive"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/apikey"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/ipaccesslist"
	serviceaccount "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/serviceaccount"
	serviceaccountaccesslistentryproject "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/serviceaccountaccesslistentry"
	serviceaccountsecret "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/project/serviceaccountsecret"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/providerconfig"
	account "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/service/account"
	accountprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/service/accountprojectassignment"
	accountsecret "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/service/accountsecret"
	partyintegration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/cluster/third/partyintegration"
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
		instance.Setup,
		provideraccessauthorization.Setup,
		provideraccesssetup.Setup,
		userorgassignment.Setup,
		userprojectassignment.Setup,
		userteamassignment.Setup,
		dnsconfigurationclusteraws.Setup,
		customrole.Setup,
		user.Setup,
		x509userauthentication.Setup,
		trigger.Setup,
		databaseinstance.Setup,
		privatelinkendpointservice.Setup,
		querylimit.Setup,
		configurationldap.Setup,
		integration.Setup,
		window.Setup,
		advancedcluster.Setup,
		auditing.Setup,
		cluster.Setup,
		cluster.Setup,
		clusterconfig.Setup,
		clusteroutagesimulation.Setup,
		employeeaccessgrant.Setup,
		organization.Setup,
		project.Setup,
		serviceaccountaccesslistentry.Setup,
		team.Setup,
		container.Setup,
		peering.Setup,
		archive.Setup,
		apikey.Setup,
		invitation.Setup,
		ipaccesslist.Setup,
		serviceaccount.Setup,
		serviceaccountaccesslistentryproject.Setup,
		serviceaccountsecret.Setup,
		providerconfig.Setup,
		account.Setup,
		accountprojectassignment.Setup,
		accountsecret.Setup,
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
		key.SetupGated,
		keyprojectassignment.SetupGated,
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
		dnsconfigurationclusteraws.SetupGated,
		customrole.SetupGated,
		user.SetupGated,
		x509userauthentication.SetupGated,
		trigger.SetupGated,
		databaseinstance.SetupGated,
		privatelinkendpointservice.SetupGated,
		querylimit.SetupGated,
		configurationldap.SetupGated,
		integration.SetupGated,
		window.SetupGated,
		advancedcluster.SetupGated,
		auditing.SetupGated,
		cluster.SetupGated,
		cluster.SetupGated,
		clusterconfig.SetupGated,
		clusteroutagesimulation.SetupGated,
		employeeaccessgrant.SetupGated,
		organization.SetupGated,
		project.SetupGated,
		serviceaccountaccesslistentry.SetupGated,
		team.SetupGated,
		container.SetupGated,
		peering.SetupGated,
		archive.SetupGated,
		apikey.SetupGated,
		invitation.SetupGated,
		ipaccesslist.SetupGated,
		serviceaccount.SetupGated,
		serviceaccountaccesslistentryproject.SetupGated,
		serviceaccountsecret.SetupGated,
		providerconfig.SetupGated,
		account.SetupGated,
		accountprojectassignment.SetupGated,
		accountsecret.SetupGated,
		partyintegration.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
