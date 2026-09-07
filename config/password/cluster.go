package password

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	v1 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

type clusterPasswordSecretRefSetter interface {
	SetPasswordSecretRef(ref *v1.SecretKeySelector)
}

// ClusterGenerator is a password initializer for cluster-scoped resources.
// Secret references include both name and namespace.
var ClusterGenerator = newGenerator(clusterSetPasswordSecretRef)

func clusterSetPasswordSecretRef(ctx context.Context, cl client.Client, mg xpresource.Managed, name, namespace, key string) error {
	setter, ok := mg.(clusterPasswordSecretRefSetter)
	if !ok {
		return nil
	}
	setter.SetPasswordSecretRef(&v1.SecretKeySelector{
		SecretReference: v1.SecretReference{Name: name, Namespace: namespace},
		Key:             key,
	})
	if err := cl.Update(ctx, mg); err != nil {
		return fmt.Errorf("cannot update managed resource with password secret ref: %w", err)
	}
	return nil
}
