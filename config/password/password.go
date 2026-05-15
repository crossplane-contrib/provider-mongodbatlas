package password

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/fieldpath"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	pw "github.com/crossplane/crossplane-runtime/v2/pkg/password"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	xpresource "github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/crossplane/upjet/v2/pkg/config"
)

type setRefFunc func(ctx context.Context, cl client.Client, mg xpresource.Managed, name, ns, key string) error

func newGenerator(setRef setRefFunc) func(string, string) config.NewInitializerFn {
	return func(byopSecretRefPath, writeConnectionSecretPath string) config.NewInitializerFn {
		return func(cl client.Client) managed.Initializer {
			return managed.InitializerFn(func(ctx context.Context, mg xpresource.Managed) error {
				paved, err := fieldpath.PaveObject(mg)
				if err != nil {
					return fmt.Errorf("cannot pave object: %w", err)
				}
				byop, err := checkBYOP(paved, byopSecretRefPath)
				if err != nil {
					return err
				}
				if byop {
					return nil
				}
				name, ns, err := resolveConnRef(paved, writeConnectionSecretPath, mg.GetNamespace())
				if err != nil || name == "" {
					return err
				}
				return reconcilePassword(ctx, cl, mg, name, ns, setRef)
			})
		}
	}
}

func checkBYOP(paved *fieldpath.Paved, path string) (bool, error) {
	sel := &v1.SecretKeySelector{}
	if err := paved.GetValueInto(path, sel); err == nil {
		return sel.Name != "", nil
	} else if xpresource.Ignore(fieldpath.IsNotFound, err) != nil {
		return false, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return false, nil
}

func resolveConnRef(paved *fieldpath.Paved, path, defaultNS string) (name, ns string, err error) {
	connRef := &v1.SecretReference{}
	if err := paved.GetValueInto(path, connRef); err != nil {
		if fieldpath.IsNotFound(err) {
			return "", "", nil
		}
		return "", "", fmt.Errorf("cannot read %s: %w", path, err)
	}
	ns = connRef.Namespace
	if ns == "" {
		ns = defaultNS
	}
	return connRef.Name, ns, nil
}

func reconcilePassword(ctx context.Context, cl client.Client, mg xpresource.Managed, name, ns string, setRef setRefFunc) error {
	const passwordKey = "password"
	s := &corev1.Secret{}
	getErr := cl.Get(ctx, types.NamespacedName{Namespace: ns, Name: name}, s)
	if xpresource.IgnoreNotFound(getErr) != nil {
		return fmt.Errorf("cannot get connection secret: %w", getErr)
	}
	if getErr == nil && len(s.Data[passwordKey]) != 0 {
		return setRef(ctx, cl, mg, name, ns, passwordKey)
	}
	return generateAndApply(ctx, cl, mg, s, name, ns, passwordKey, setRef)
}

func generateAndApply(ctx context.Context, cl client.Client, mg xpresource.Managed, s *corev1.Secret, name, ns, key string, setRef setRefFunc) error {
	generated, err := pw.Generate()
	if err != nil {
		return fmt.Errorf("cannot generate password: %w", err)
	}
	s.SetName(name)
	s.SetNamespace(ns)
	s.Type = xpresource.SecretTypeConnection
	if !meta.WasCreated(s) {
		meta.AddOwnerReference(s, meta.AsController(meta.TypedReferenceTo(mg, mg.GetObjectKind().GroupVersionKind())))
	}
	if s.Data == nil {
		s.Data = make(map[string][]byte, 1)
	}
	s.Data[key] = []byte(generated)
	if err := xpresource.NewAPIPatchingApplicator(cl).Apply(ctx, s); err != nil {
		return fmt.Errorf("cannot apply password secret: %w", err)
	}
	return setRef(ctx, cl, mg, name, ns, key)
}
