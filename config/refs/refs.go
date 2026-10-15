package refs

const (
	APIKeyID       = "api_key_id"
	ClusterName    = "cluster_name"
	ConnectionName = "connection_name"
	ContainerID    = "container_id"
	IAMRoleID      = "iam_role_id"
	InstanceName   = "instance_name"
	Name           = "name"
	OrgID          = "org_id"
	PeerID         = "peer_id"
	ProjectID      = "project_id"
	ClientID       = "client_id"
	ProviderName   = "provider_name"
	Region         = "region"
	RoleName       = "role_name"
	SnapshotID     = "snapshot_id"
	TeamID         = "team_id"
	WorkspaceName  = "workspace_name"

	TFBackupSnapshot              = "mongodbatlas_cloud_backup_snapshot"
	TFCluster                     = "mongodbatlas_advanced_cluster"
	TFOrganization                = "mongodbatlas_organization"
	TFProject                     = "mongodbatlas_project"
	TFProjectServiceAccount       = "mongodbatlas_project_service_account"
	TFProviderAccessAuthorization = "mongodbatlas_cloud_provider_access_authorization"
	TFServiceAccount              = "mongodbatlas_service_account"
	TFStreamConnection            = "mongodbatlas_stream_connection"
	TFStreamInstance              = "mongodbatlas_stream_instance"
	TFStreamWorkspace             = "mongodbatlas_stream_workspace"
	TFTeam                        = "mongodbatlas_team"

	VersionV1Alpha2 = "v1alpha2"
	VersionV1Alpha3 = "v1alpha3"

	ExtractParamPathFmt = `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath(%q,%t)`
)
