/*
Copyright 2021 Upbound Inc.
*/

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/upbound/upjet/pkg/controller"

	listapikey "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/access/listapikey"
	configuration "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/alert/configuration"
	key "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/api/key"
	compliancepolicy "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/backup/compliancepolicy"
	backupschedule "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupschedule"
	backupsnapshot "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupsnapshot"
	backupsnapshotexportbucket "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupsnapshotexportbucket"
	backupsnapshotexportjob "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/cloud/backupsnapshotexportjob"
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
	settingsidentityprovider "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/federated/settingsidentityprovider"
	settingsorgconfig "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/federated/settingsorgconfig"
	settingsorgrolemapping "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/federated/settingsorgrolemapping"
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
	endpointregionalmode "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/private/endpointregionalmode"
	ipmode "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/private/ipmode"
	endpoint "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpoint"
	endpointserverless "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpointserverless"
	endpointservice "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpointservice"
	endpointserviceadl "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpointserviceadl"
	endpointserviceserverless "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/privatelink/endpointserviceserverless"
	apikey "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/project/apikey"
	invitationproject "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/project/invitation"
	ipaccesslist "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/project/ipaccesslist"
	providerconfig "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/providerconfig"
	index "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/search/index"
	instance "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/serverless/instance"
	partyintegration "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/third/partyintegration"
	authenticationdatabaseuser "github.com/crossplane-contrib/provider-jet-mongodbatlas/internal/controller/x509/authenticationdatabaseuser"
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
