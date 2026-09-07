package password

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	v2 "github.com/crossplane/crossplane/apis/v2/core/v2"
)

type namespacedPasswordSecretRefSetter interface {
	SetPasswordSecretRef(ref *v2.LocalSecretKeySelector)
}

// NamespacedGenerator is a password initializer for namespaced resources.
// Secret references include only name (same namespace as the CR).
var NamespacedGenerator = newGenerator(namespacedSetPasswordSecretRef)

func namespacedSetPasswordSecretRef(ctx context.Context, cl client.Client, mg xpresource.Managed, name, _, key string) error {
	setter, ok := mg.(namespacedPasswordSecretRefSetter)
	if !ok {
		return nil
	}
	setter.SetPasswordSecretRef(&v2.LocalSecretKeySelector{
		LocalSecretReference: v2.LocalSecretReference{Name: name},
		Key:                  key,
	})
	if err := cl.Update(ctx, mg); err != nil {
		return fmt.Errorf("cannot update managed resource with password secret ref: %w", err)
	}
	return nil
}
