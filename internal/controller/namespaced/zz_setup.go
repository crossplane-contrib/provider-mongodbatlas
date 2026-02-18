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
	clusterconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/global/clusterconfig"
	configurationldap "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/ldap/configuration"
	integration "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/log/integration"
	window "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/maintenance/window"
	advancedcluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/advancedcluster"
	auditing "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/auditing"
	cluster "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/cluster"
	employeeaccessgrant "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/employeeaccessgrant"
	organization "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/organization"
	project "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/project"
	team "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/mongodbatlas/team"
	container "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/container"
	peering "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/network/peering"
	archive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/online/archive"
	endpointservicedatafederationonlinearchive "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/privatelink/endpointservicedatafederationonlinearchive"
	apikey "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/apikey"
	invitation "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/ipaccesslist"
	serviceaccount "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/serviceaccount"
	serviceaccountaccesslistentry "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/serviceaccountaccesslistentry"
	serviceaccountsecret "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/project/serviceaccountsecret"
	providerconfig "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/providerconfig"
	instance "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/serverless/instance"
	account "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/service/account"
	accountaccesslistentry "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/service/accountaccesslistentry"
	accountprojectassignment "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/service/accountprojectassignment"
	accountsecret "github.com/crossplane-contrib/provider-mongodbatlas/internal/controller/namespaced/service/accountsecret"
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
		clusterconfig.Setup,
		configurationldap.Setup,
		integration.Setup,
		window.Setup,
		advancedcluster.Setup,
		auditing.Setup,
		cluster.Setup,
		cluster.Setup,
		employeeaccessgrant.Setup,
		organization.Setup,
		project.Setup,
		team.Setup,
		container.Setup,
		peering.Setup,
		archive.Setup,
		endpointservicedatafederationonlinearchive.Setup,
		apikey.Setup,
		invitation.Setup,
		ipaccesslist.Setup,
		serviceaccount.Setup,
		serviceaccountaccesslistentry.Setup,
		serviceaccountsecret.Setup,
		providerconfig.Setup,
		instance.Setup,
		account.Setup,
		accountaccesslistentry.Setup,
		accountprojectassignment.Setup,
		accountsecret.Setup,
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
		clusterconfig.SetupGated,
		configurationldap.SetupGated,
		integration.SetupGated,
		window.SetupGated,
		advancedcluster.SetupGated,
		auditing.SetupGated,
		cluster.SetupGated,
		cluster.SetupGated,
		employeeaccessgrant.SetupGated,
		organization.SetupGated,
		project.SetupGated,
		team.SetupGated,
		container.SetupGated,
		peering.SetupGated,
		archive.SetupGated,
		endpointservicedatafederationonlinearchive.SetupGated,
		apikey.SetupGated,
		invitation.SetupGated,
		ipaccesslist.SetupGated,
		serviceaccount.SetupGated,
		serviceaccountaccesslistentry.SetupGated,
		serviceaccountsecret.SetupGated,
		providerconfig.SetupGated,
		instance.SetupGated,
		account.SetupGated,
		accountaccesslistentry.SetupGated,
		accountprojectassignment.SetupGated,
		accountsecret.SetupGated,
		partyintegration.SetupGated,
		authenticationdatabaseuser.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
