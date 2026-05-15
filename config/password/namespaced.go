package password

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
)

type namespacedPasswordSecretRefSetter interface {
	SetPasswordSecretRef(ref *v1.LocalSecretKeySelector)
}

// NamespacedGenerator is a password initializer for namespaced resources.
// Secret references include only name (same namespace as the CR).
var NamespacedGenerator = newGenerator(namespacedSetPasswordSecretRef)

func namespacedSetPasswordSecretRef(ctx context.Context, cl client.Client, mg xpresource.Managed, name, _, key string) error {
	setter, ok := mg.(namespacedPasswordSecretRefSetter)
	if !ok {
		return nil
	}
	setter.SetPasswordSecretRef(&v1.LocalSecretKeySelector{
		LocalSecretReference: v1.LocalSecretReference{Name: name},
		Key:                  key,
	})
	if err := cl.Update(ctx, mg); err != nil {
		return fmt.Errorf("cannot update managed resource with password secret ref: %w", err)
	}
	return nil
}
