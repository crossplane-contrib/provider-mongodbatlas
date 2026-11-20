package resources

import (
	"testing"

	"github.com/crossplane/upjet/v2/pkg/resource"
	"github.com/stretchr/testify/assert"

	"github.com/crossplane-contrib/provider-mongodbatlas/apis/cluster/mongodbatlas/v1alpha3"
)

// TestFirstReplicationSpecZoneID asserts that the extractor of the BackupSchedule copySettings zoneId reference
// reads the zone of the referenced AdvancedCluster's first replication spec from its status.
func TestFirstReplicationSpecZoneID(t *testing.T) {
	zone := func(id string) v1alpha3.ReplicationSpecsObservation {
		return v1alpha3.ReplicationSpecsObservation{ZoneID: &id}
	}
	cluster := &v1alpha3.AdvancedCluster{}
	cluster.Status.AtProvider.ReplicationSpecs = []v1alpha3.ReplicationSpecsObservation{zone("first"), zone("second")}

	assert.Equal(t, "first", resource.ExtractParamPath(firstReplicationSpecZoneID, true)(cluster))
	assert.Empty(t, resource.ExtractParamPath(firstReplicationSpecZoneID, true)(&v1alpha3.AdvancedCluster{}),
		"a cluster without an observed replication spec must give no zone")
}
